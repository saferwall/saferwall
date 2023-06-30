// Copyright 2018 Saferwall. All rights reserved.
// Use of this source code is governed by Apache v2 license
// license that can be found in the LICENSE file.

package sandbox

import (
	"bytes"
	"context"
	"image"
	"image/color"
	"io"
	"path/filepath"
	"strconv"

	"github.com/disintegration/imaging"
	agent "github.com/saferwall/saferwall/internal/agent"
	"github.com/saferwall/saferwall/internal/log"
	"github.com/saferwall/saferwall/internal/utils"
	"github.com/saferwall/saferwall/services/config"
)

const (
	defaultGRPCPort        = ":50051"
	defaultFileScanTimeout = 30
	defaultVPNCountry      = "USA"
	defaultOS              = "Windows 7 64-bit"
)

// EventType is the type of the system event. A type can be either:
// `registry`, `network` or `file`.
type EventType string

const (
	FileEventType     = "file"
	RegistryEventType = "registry"
	NetworkEventType  = "network"
)

// DetonationResult represents the results for a detonation.
type DetonationResult struct {
	// The API trace results. This consists of a list of all API calls made by
	// the sample.
	APITrace []byte `json:"api_trace,omitempty"`
	// The logs produced by the agent running inside the VM.
	AgentLog []byte `json:"agent_log,omitempty"`
	// The logs produced by the sandbox.
	SandboxLog []byte `json:"sandbox_log,omitempty"`
	// The config used to scan dynamically the sample.
	ScanCfg config.DynFileScanCfg `json:"scan_config,omitempty"`
	// List of of desktop screenshots captured.
	Screenshots []Screenshot `json:"screenshots,omitempty"`
	// List of artifacts collected during detonation.
	Artifacts []Artifact `json:"artifacts,omitempty"`
	// Environment represents the environment used to scan the file.
	Environment `json:"env,omitempty"`
	// Summary of system events.
	Events map[EventType][]Event `json:"events"`
}

// Environment represents the environment used to scan the file.
// This include OS versions, installed software, analyzer version.
type Environment struct {
	// The sandbox version.
	SandboxVersion string `json:"sandbox_version"`
	// The agent version.
	AgentVersion string `json:"agent_version"`
}

// Screenshot represents a capture of the desktop while the sample is running
// in the VM.
type Screenshot struct {
	// The name of the filename for the screenshot. Format: <id>.jpeg.
	// IDs are growing incrementally from index 1 to N according to the time
	// they were taken.
	Name string `json:"name"`
	// The binary content of the image.
	Content []byte `json:"-"`
}

// Artifact represents dumped memory buffers (during process injection, memory
// decryption, or anything alike) and files dropped by the sample.
type Artifact struct {
	// File  name of the artifact.
	// * Memory buffers are in this format: PID.TID-VA-BuffSize-API.membuff
	//  -> 2E00.A60-0x1A46D880000-77824-Free.membuff
	//  -> 2E00.A60-0x1A46D8A0000-12824-Crypto.membuff
	//  -> 2E00.A60-0x1A46D9B0000-12824-WriteProcess.membuff
	// * Files dropped are in this format: PID.TID-FilePath-Size.filecreate
	//  -> 2E00.A60-C:\\Delete.vbs-9855.filecreate
	// they were taken.
	Name string `json:"name"`
	// The binary content of the artifact.
	Content []byte `json:"-"`
	// The artifact kind.
	Kind string `json:"kind"`
	// The SHA256 hash of the artifact.
	SHA256 string `json:"sha256"`
	// Detection contains the family name of the malware if it is malicious,
	// or clean otherwise.
	Detection string `json:"detection"`
	// The file type, i.e docx, dll, etc.
	FileType string `json:"file_type"`
}

// Event represents a system event: a registry, network or file event.
type Event struct {
	// Path of the system event. For instance, when the event is of type:
	// `registry`, the path represents the registry key being used. For a
	// `network` event type, the path is the IP or domain used.
	Path string `json:"path"`
	// Th operation requested over the above `Path` field. This field means
	// different things according to the type of the system event.
	// - For file system events: can be either: create, read, delete, rename,
	// - For registry events: can be either: create, rename, set, delete.
	// - For network events: can be either an HTTP verb is the network is based
	// of HTTP or connect, send, receive if communication happens over plain sockets.
	Operation string `json:"op"`
}

func (s *Service) detonate(logger log.Logger, vm *VM,
	cfg config.FileScanCfg) (DetonationResult, error) {

	detRes := DetonationResult{}
	ctx := context.Background()

	// Establish a gRPC connection to the agent server running inside the guest.
	client, err := agent.New(vm.IP + defaultGRPCPort)
	if err != nil {
		logger.Errorf("failed to establish connection to server: %v", err)
		return detRes, err
	}

	// Deploy the sandbox component files inside the guest.
	ver, err := client.Deploy(ctx, s.cfg.Agent.AgentDestDir, s.sandbox)
	if err != nil {
		return detRes, err
	}
	logger.Infof("sandbox version %s has been deployed", ver)

	detRes.SandboxVersion = ver
	detRes.AgentVersion = vm.AgentVersion

	src := filepath.Join(s.cfg.SharedVolume, cfg.SHA256)
	sampleContent, err := utils.ReadAll(src)
	if err != nil {
		return detRes, err
	}

	// Analyze the sample. This call will block until results are ready.
	scanCfg := toJSON(cfg.DynFileScanCfg)
	detRes.ScanCfg = cfg.DynFileScanCfg
	res, err := client.Analyze(ctx, scanCfg, sampleContent)
	if err != nil {
		return detRes, err
	}

	// Convert the agent log from JSONL to JSON.
	var agentLog []interface{}
	err = Decode(bytes.NewReader(res.AgentLog), &agentLog)
	if err != nil {
		logger.Errorf("failed to decode agent log: %v", err)
	}
	detRes.AgentLog = toJSON(agentLog)

	// Convert the APIs traces from JSONL to JSON.
	var traceLog []Win32API
	err = Decode(bytes.NewReader(res.TraceLog), &traceLog)
	if err != nil {
		logger.Errorf("failed to decode trace log: %v", err)
	}
	detRes.APITrace = toJSON(traceLog)

	// Convert the sandbox log from JSONL to JSON.
	var sandboxLog []interface{}
	err = Decode(bytes.NewReader(res.SandboxLog), &sandboxLog)
	if err != nil {
		logger.Errorf("failed to decode sandbox log: %v", err)
	}
	detRes.SandboxLog = toJSON(sandboxLog)

	// Create a summary of system events.
	detRes.Events, err = s.summarizeEvents(traceLog)
	if err != nil {
		logger.Errorf("failed to summarize behavior events: %v", err)
	}

	// Collect screenshots.
	screenshots := []Screenshot{}
	for _, sc := range res.Screenshots {
		screenshots = append(screenshots, Screenshot{
			Name:    strconv.Itoa(int(sc.GetId())) + ".jpeg",
			Content: sc.GetContent(),
		})

		// Generate thumbnails.
		r := bytes.NewReader(sc.GetContent())
		buf := new(bytes.Buffer)
		err := s.generateThumbnail(r, buf)
		if err != nil {
			logger.Errorf("failed to generate thumbnail: %v", err)
			continue
		}

		screenshots = append(screenshots, Screenshot{
			Name:    strconv.Itoa(int(sc.GetId())) + ".min.jpeg",
			Content: buf.Bytes(),
		})
	}
	detRes.Screenshots = screenshots

	// Collect artifacts.
	artifacts := []Artifact{}
	for _, artifact := range res.MemDumps {
		artifacts = append(artifacts, Artifact{
			Name:    artifact.GetName(),
			Content: artifact.GetContent(),
		})
	}
	detRes.Artifacts = artifacts

	return detRes, nil
}

// generate thumbnails for the sandbox desktop screenshots.
func (s *Service) generateThumbnail(r io.Reader, w io.Writer) error {

	// load images and make 100x100 thumbnails of them.
	img, err := imaging.Decode(r, nil...)
	if err != nil {
		return err
	}

	x := 730
	y := 450
	thumbnail := imaging.Thumbnail(img, x, y, imaging.CatmullRom)

	// create a new blank image
	dst := imaging.New(x, y, color.NRGBA{0, 0, 0, 0})

	// paste thumbnails into the new image side by side
	dst = imaging.Paste(dst, thumbnail, image.Pt(0, 0))

	// write the combined image to an io writer.
	opts := imaging.JPEGQuality(80)
	err = imaging.Encode(w, dst, imaging.JPEG, opts)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) summarizeEvents(w32apis []Win32API) (
	map[EventType][]Event, error) {

	events := make(map[EventType][]Event)

	for _, w32api := range w32apis {
		// Registry events
		if utils.StringInSlice(w32api.Name, regAPIs) {
			event := summarizeRegAPI(w32api)
			events[RegistryEventType] = append(events[RegistryEventType], event)
		}
	}

	return events, nil
}
