# 📂 GoFiles

![Go Version](https://img.shields.io/badge/Go-1.25.5-blue.svg) ![License](https://img.shields.io/badge/License-MIT-green.svg)

**GoFiles** is a robust, lightweight, and modern file management server built with **Go** (Golang). It provides a blazing-fast RESTful API to browse, search, manage, and organize your files with ease. Featuring a built-in safe trash system, content searching, and secure file operations, it's designed to be the backend for your next file manager or personal cloud.

---

## ✨ Key Features

- **🚀 High Performance:** Built on Go's standard library for minimal footprint and maximum speed.
- **📂 File Browsing:** List files and directories with rich metadata (size, modification time, type).
- **🔍 Advanced Search:**
  - **Name Search:** Find files instantly by name.
  - **Content Search:** Deep search within text files (skips large files > 5MB for performance).
- **🛡️ Trash System:** Safe deletion with a dedicated `.trash` folder.
  - Soft delete (move to trash).
  - Restore files.
  - Auto-expiration (default 30 days).
- **🛠️ Full Management:**
  - **Upload** files.
  - **Create** directories.
  - **Rename**, **Move**, and **Copy** files/folders.
  - **Download** files securely.
- **⚙️ Configurable:** safe root directory confinement.

---

## 🛠️ Installation & Usage

### Prerequisites

- [Go](https://go.dev/dl/) installed (v1.25+ recommended).

### Running Locally

1.  Clone the repository:

    ```bash
    git clone https://github.com/yourusername/GoFiles.git
    cd GoFiles
    ```

2.  Run the server:

    ```bash
    go run .
    ```

3.  The server will start on port **8080**:
    ```
    🚀 GoFiles Server started on http://localhost:8080
    ```

---

## ⚙️ Configuration

Configuration is currently handled via constants in standard `config.go`.

| Constant         | Default   | Description                                                                        |
| :--------------- | :-------- | :--------------------------------------------------------------------------------- |
| `RootFolder`     | `.`       | The root directory to serve files from. Defaults to the current running directory. |
| `TrashFolder`    | `.trash`  | The hidden directory used for storing deleted files.                               |
| `TrashRetention` | `30 Days` | Duration before trashed files are considered for permanent removal.                |

---

## 📖 API Documentation

The server exposes a comprehensive REST API. All responses are in JSON format.

### 🔍 Read & Search

| Method | Endpoint        | Query Params                                                   | Description                          |
| :----- | :-------------- | :------------------------------------------------------------- | :----------------------------------- |
| `GET`  | `/api/files`    | `path` (relative), `ext` (filter), `min_size` (bytes)          | List files in a directory.           |
| `GET`  | `/api/search`   | `q` (query), `type` (`name` or `content`), `path` (start path) | Search for files by name or content. |
| `GET`  | `/api/download` | `path`                                                         | Download a specific file.            |

### ✍️ Write & Upload

| Method   | Endpoint      | Body / Form                              | Description                                           |
| :------- | :------------ | :--------------------------------------- | :---------------------------------------------------- |
| `POST`   | `/api/upload` | Form-Data: `file`                        | Upload a file to the directory specified by `?path=`. |
| `POST`   | `/api/mkdir`  | JSON: `{ "path": "...", "name": "..." }` | Create a new directory.                               |
| `DELETE` | `/api/delete` | Query: `path`, `permanent=true/false`    | Delete a file/folder. Defaults to moving to trash.    |

### 📦 Organize

| Method | Endpoint      | Body (JSON)                                  | Description                 |
| :----- | :------------ | :------------------------------------------- | :-------------------------- |
| `POST` | `/api/rename` | `{ "sourcePath": "...", "newName": "..." }`  | Rename a file or directory. |
| `POST` | `/api/move`   | `{ "sourcePath": "...", "destPath": "..." }` | Move a file or directory.   |
| `POST` | `/api/copy`   | `{ "sourcePath": "...", "destPath": "..." }` | Copy a file or directory.   |

### 🗑️ Trash Management

| Method | Endpoint             | Query Params            | Description                                         |
| :----- | :------------------- | :---------------------- | :-------------------------------------------------- |
| `GET`  | `/api/trash/list`    | -                       | List all files currently in the trash.              |
| `POST` | `/api/trash/restore` | `name` (trash filename) | Restore a file from trash to its original location. |
| `POST` | `/api/trash/empty`   | -                       | Permanently delete all files in the trash.          |

---

## 📂 Project Structure

```bash
GoFiles/
├── config.go          # Configuration constants (Root folder, Trash settings)
├── handlers_ops.go    # Logic for Move, Copy, Rename, and Delete
├── handlers_read.go   # Logic for Listing, Searching, and Downloading
├── handlers_trash.go  # Logic for Trash bin management (List, Restore, Empty)
├── handlers_write.go  # Logic for Uploads and Directory creation
├── main.go            # Entry point & Router setup
├── trash.go           # Internal trash utilities (MoveToTrash, RestoreFromTrash)
├── types.go           # Struct definitions (API Requests/Responses)
├── utils.go           # Helper functions (CORS, Safe Path checks)
└── go.mod             # Go module definition
```

---

## 📄 License

This project is licensed under the [MIT License](LICENSE).

---

Made with ❤️ in Go.
