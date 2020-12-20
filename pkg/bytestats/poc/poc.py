import numpy as np


def byte_histogram(bytez):
    counts = np.bincount(np.frombuffer(bytez, dtype=np.uint8), minlength=256)
    return counts.tolist()


def _entropy_bin_counts(window, block):
    # coarse histogram, 16 bytes per bin
    c = np.bincount(block >> 4, minlength=16)  # 16-bin histogram
    p = c.astype(np.float32) / window
    wh = np.where(c)[0]
    H = (
        np.sum(-p[wh] * np.log2(p[wh])) * 2
    )  # * x2 b.c. we reduced information by half: 256 bins (8 bits) to 16 bins (4 bits)

    Hbin = int(H * 2)  # up to 16 bins (max entropy is 8 bits)
    if Hbin == 16:  # handle entropy = 8.0 bits
        Hbin = 15

    return Hbin, c


def byte_entropy_histogram(bytez):
    step = 1024
    window = 2048
    output = np.zeros((16, 16), dtype=np.int)
    a = np.frombuffer(bytez, dtype=np.uint8)
    if a.shape[0] < window:
        Hbin, c = _entropy_bin_counts(window, a)
        output[Hbin, :] += c
    else:
        # strided trick from here: http://www.rigtorp.se/2011/01/01/rolling-statistics-numpy.html
        shape = a.shape[:-1] + (a.shape[-1] - window + 1, window)
        strides = a.strides + (a.strides[-1],)
        blocks = np.lib.stride_tricks.as_strided(
            a, shape=shape, strides=strides
        )[::step, :]

        # from the blocks, compute histogram
        for block in blocks:
            Hbin, c = _entropy_bin_counts(window, block)
            output[Hbin, :] += c

    return output.flatten().tolist()


if __name__ == "__main__":
    path = r"C:\Users\kaplan\Projects\saferwall\binaries\cmd.exe"
    bytez = np.memmap(path)
    print(byte_entropy_histogram(bytez))