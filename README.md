# Media Dedup CLI

Imagine you having your iPhone Photo library full of duplicates. How to clean it? Here is what this tool needed.

A high-performance **Go** command-line tool designed for macOS to identify and remove visual duplicate photos (including Apple HEIC) and identical video files.

### Key Features
* **Perceptual Hashing (pHash):** Finds duplicate photos even if they have been resized, compressed, or renamed.
* **Native HEIC Support:** Seamlessly processes iPhone photos using the macOS system utility `sips`.
* **Fast Video Fingerprinting:** Instantly matches large video files (MP4/MOV) by comparing file size and header/footer signatures without reading the entire file.
* **Safety First:** Runs in "Dry Run" mode by default. Deletion only occurs when explicitly authorized.

---

## 🛠 Prerequisites
* **OS:** macOS (required for `sips` integration in case if you have HEIC files).
* **Language:** Go 1.25 or higher.

---

## 🚀 Installation & Setup

1.  **Initialize the module:**
    ```bash
    go mod init media-dedup
    go mod tidy
    ```

2.  **Build the binary:**
    ```bash
    go build -o dedup main.go
    ```

---

## 💻 Usage

The tool uses flags to define the scan scope and action.

### 1. Scan a Specific Directory (Dry Run)
This will list all duplicates found without deleting anything.
```bash
./dedup -dir /Users/vm/Pictures/iphone
```

### 2. Scan and DELETE Duplicates
Use the `-delete` flag to trigger permanent removal of the identified duplicates.
```bash
./dedup -dir /Users/vm/Pictures/iphone -delete
```

### 3. Scan Current Directory
If no directory is provided, it defaults to the folder where the binary is executed.
```bash
./dedup -delete
```

---

## ⚙️ Arguments

| Flag | Description | Default |
| :--- | :--- | :--- |
| `-dir` | The absolute or relative path to the folder you want to scan. | `.` (Current) |
| `-delete` | Enable permanent deletion of found duplicates. | `false` |

---

## 🧠 How It Works

### For Photos (`.jpg`, `.png`, `.heic`)
The tool generates a **64-bit Perceptual Hash**. Unlike standard MD5 hashes, pHash focuses on the visual structure. If you have the same photo saved from iCloud and a Telegram cache, they will be matched regardless of metadata or filename differences.



### For Videos (`.mp4`, `.mov`, `.m4v`)
To avoid bottlenecking your SSD by reading multi-gigabyte files, the tool uses a **Fast Signature** method:
1.  Matches exact **File Size**.
2.  Samples and hashes the **first 100KB** and **last 100KB** of the file.
This is 99.9% accurate for identifying duplicates like `MOV_01.mp4` vs `MOV_01 1.mp4`.

---

> **Warning:** Permanent deletion cannot be undone. Always run the tool without the `-delete` flag first to verify the results.
