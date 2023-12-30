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
	defaultFileScanTimeout = 10
	defaultVPNCountry      = "USA"
	defaultOS              = "Windows 7 64-bit"
	maxTraceLog            = 10000
)

const (
	fileEventType     = "file"
	registryEventType = "registry"
	networkEventType  = "network"
)

// DetonationResult represents the results for a detonation.
type DetonationResult struct {
	// The API trace results. This consists of a list of all API calls made by
	// the sample.
	APITrace []byte
	// Same as APITrace, but this is not capped to any threshold.
	// The full trace is uploaded to the object storage,
	FullAPITrace []byte
	// The buffer of large byte* for some Win32 APIs.
	APIBuffers []Win32APIBuffer
	// The logs produced by the agent running inside the VM.
	AgentLog []byte
	// The logs produced by the sandbox.
	SandboxLog []byte
	// The config used to scan dynamically the sample.
	ScanCfg config.DynFileScanCfg
	// List of desktop screenshots captured.
	Screenshots []Screenshot
	// List of the sample capabilities.
	Capabilities []Capability
	// List of artifacts collected during detonation.
	Artifacts []Artifact
	// Environment represents the environment used to scan the file.
	Environment
	// Summary of system events.
	Events []Event
	// The process execution tree.
	ProcessTree ProcessTree
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
	Name string
	// The binary content of the image.
	Content []byte
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

	// TODO: After the analysis is finish, we can immediately start reverting
	// the snapshot in a goroutine and get the result later.

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

	// Collect API buffers.
	for _, apiBuff := range res.APIBuffers {
		detRes.APIBuffers = append(detRes.APIBuffers, Win32APIBuffer{
			Name:    apiBuff.GetName(),
			Content: apiBuff.GetContent(),
		})
	}

	// Convert the sandbox log from JSONL to JSON.
	var sandboxLog []interface{}
	err = Decode(bytes.NewReader(res.SandboxLog), &sandboxLog)
	if err != nil {
		logger.Errorf("failed to decode sandbox log: %v", err)
	}
	detRes.SandboxLog = toJSON(sandboxLog)

	// Create the process tree structure.
	var processes []Process
	err = Decode(bytes.NewReader(res.ProcessTree), &processes)
	if err != nil {
		logger.Errorf("failed to decode process tree: %v", err)
	}

	// TODO: Detect API calls in loops ! The JSON log is capped to 20MB.
	var minTraceLog []Win32API
	detRes.FullAPITrace = toJSON(traceLog)
	if len(traceLog) > maxTraceLog {
		minTraceLog = traceLog[:maxTraceLog]
	}

	// Create a optimized version of the API trace for storage in DB.
	detRes.APITrace = curateAPIEvents(minTraceLog)

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

	// Run behavior rule scan & extract system events.
	bhvRulesMatch, err := s.bhvScanner.Scan(detRes.FullAPITrace)
	if err != nil {
		logger.Errorf("failed to scan with behavior with: %v", err)
	}
	detRes.Events = generateEvents(bhvRulesMatch.Events)

	// Run yara scan over the artifacts.
	yaraRulesMatch, err := s.scanArtifactsWithYara(res.Artifacts)
	if err != nil {
		logger.Errorf("failed to scan with artifacts with yara: %v", err)
	}
	// Generate artifacts.
	detRes.Artifacts, err = s.generateArtifacts(res.Artifacts)
	if err != nil {
		logger.Errorf("failed to generate artifacts metadata: %v", err)
	}

	// Generate capabilities.
	detRes.Capabilities = generateCapabilities(bhvRulesMatch.Rules, yaraRulesMatch)

	// This has to run last, by the time all rules / behavior are extarcted.
	detRes.ProcessTree = enrichProcTree(processes, detRes.Capabilities)

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
