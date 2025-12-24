/* eslint-disable react-hooks/exhaustive-deps */
import { useEffect, useState } from "react";
import { api, type FileInfo } from "./api";
import FileIcon from "./FileIcon";
import {
  Search,
  Home,
  ArrowUp,
  RefreshCw,
  LogOut,
  Download,
  Trash2,
  FolderPlus,
  FilePlus,
  Grid,
  List as ListIcon,
  Folder,
} from "lucide-react";
import { formatDistanceToNow } from "date-fns";

export default function Dashboard() {
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [currentPath, setCurrentPath] = useState("");
  const [loading, setLoading] = useState(false);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [searchQuery, setSearchQuery] = useState("");

  // Fetch files when path changes
  useEffect(() => {
    loadFiles();
  }, [currentPath]);

  const loadFiles = async () => {
    setLoading(true);
    try {
      const res = await api.get(`/files?path=${currentPath}`);
      setFiles(res.data);
    } catch (err) {
      console.error("Failed to load files", err);
    } finally {
      setLoading(false);
    }
  };

  const handleNavigate = (folderName: string) => {
    setCurrentPath((prev) => (prev ? `${prev}/${folderName}` : folderName));
  };

  const handleGoUp = () => {
    if (!currentPath) return;
    const parts = currentPath.split("/");
    parts.pop();
    setCurrentPath(parts.join("/"));
  };

  const handleDownload = (fileName: string) => {
    // We use window.open to trigger the browser's download behavior
    // Using the backend API URL directly
    const downloadUrl = `${api.defaults.baseURL}/download?path=${
      currentPath ? currentPath + "/" : ""
    }${fileName}`;
    window.open(downloadUrl, "_blank");
  };

  return (
    <div className="min-h-screen bg-gray-950 flex flex-col">
      {/* --- HEADER --- */}
      <header className="bg-gray-900 border-b border-gray-800 px-6 py-4 flex items-center justify-between sticky top-0 z-10">
        <div className="flex items-center gap-4">
          <div className="flex items-center gap-2 text-xl font-bold text-white">
            <span className="bg-blue-600 text-transparent bg-clip-text">
              GoFiles
            </span>
          </div>

          {/* Breadcrumbs */}
          <div className="hidden md:flex items-center gap-2 ml-8 px-4 py-2 bg-gray-950 rounded-lg border border-gray-800 text-sm text-gray-400">
            <button
              onClick={() => setCurrentPath("")}
              className="hover:text-white transition-colors"
            >
              <Home className="w-4 h-4" />
            </button>
            <span className="text-gray-700">/</span>
            {currentPath ? (
              <span className="text-white truncate max-w-[200px]">
                {currentPath}
              </span>
            ) : (
              <span>Home</span>
            )}
          </div>
        </div>

        <div className="flex items-center gap-3">
          {/* Search Bar */}
          <div className="relative hidden sm:block">
            <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-500" />
            <input
              type="text"
              placeholder="Search..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="bg-gray-950 border border-gray-800 text-sm rounded-full pl-10 pr-4 py-2 text-gray-300 focus:outline-none focus:border-blue-500 w-64 transition-all"
            />
          </div>

          <div className="h-6 w-px bg-gray-800 mx-2"></div>

          <button
            onClick={() => setViewMode(viewMode === "grid" ? "list" : "grid")}
            className="p-2 text-gray-400 hover:text-white hover:bg-gray-800 rounded-lg transition-all"
          >
            {viewMode === "grid" ? (
              <ListIcon className="w-5 h-5" />
            ) : (
              <Grid className="w-5 h-5" />
            )}
          </button>

          <button
            onClick={() => window.location.reload()}
            className="p-2 text-red-400 hover:bg-red-500/10 rounded-lg transition-all"
          >
            <LogOut className="w-5 h-5" />
          </button>
        </div>
      </header>

      {/* --- TOOLBAR --- */}
      <div className="px-6 py-3 border-b border-gray-800 flex items-center gap-2 overflow-x-auto">
        <button
          onClick={handleGoUp}
          disabled={!currentPath}
          className="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-gray-300 hover:bg-gray-800 rounded-md disabled:opacity-50 disabled:cursor-not-allowed"
        >
          <ArrowUp className="w-4 h-4" /> Up
        </button>
        <button
          onClick={loadFiles}
          className="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-gray-300 hover:bg-gray-800 rounded-md"
        >
          <RefreshCw className={`w-4 h-4 ${loading ? "animate-spin" : ""}`} />{" "}
          Refresh
        </button>
        <div className="w-px h-4 bg-gray-800 mx-2"></div>
        <button className="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-blue-400 hover:bg-blue-500/10 rounded-md">
          <FolderPlus className="w-4 h-4" /> New Folder
        </button>
        <button className="flex items-center gap-2 px-3 py-1.5 text-sm font-medium text-green-400 hover:bg-green-500/10 rounded-md">
          <FilePlus className="w-4 h-4" /> Upload
        </button>
      </div>

      {/* --- MAIN CONTENT --- */}
      <main className="flex-1 p-6 overflow-y-auto">
        {files.length === 0 && !loading ? (
          <div className="flex flex-col items-center justify-center h-64 text-gray-500">
            <Folder className="w-16 h-16 mb-4 text-gray-700" />
            <p>This folder is empty</p>
          </div>
        ) : (
          <div
            className={
              viewMode === "grid"
                ? "grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4"
                : "flex flex-col gap-2"
            }
          >
            {files.map((file) => (
              <div
                key={file.name}
                onClick={() => file.is_dir && handleNavigate(file.name)}
                className={`
                                    group relative border border-gray-800/50 hover:border-blue-500/50 hover:bg-gray-900 rounded-xl transition-all cursor-pointer select-none
                                    ${
                                      viewMode === "grid"
                                        ? "p-4 flex flex-col items-center text-center aspect-square justify-center"
                                        : "p-3 flex items-center justify-between"
                                    }
                                `}
              >
                <div
                  className={`flex items-center ${
                    viewMode === "grid" ? "flex-col gap-3" : "gap-4"
                  }`}
                >
                  {/* Thumbnail Logic: If image, try to load thumb, else show icon */}
                  {["jpg", "jpeg", "png", "gif"].includes(
                    file.type.replace(".", "").toLowerCase()
                  ) ? (
                    <img
                      src={`${api.defaults.baseURL}/thumbnail?path=${
                        currentPath ? currentPath + "/" : ""
                      }${file.name}`}
                      alt={file.name}
                      className={`${
                        viewMode === "grid"
                          ? "w-16 h-16 object-cover rounded-lg"
                          : "w-10 h-10 rounded object-cover"
                      }`}
                      loading="lazy"
                    />
                  ) : (
                    <FileIcon
                      isDir={file.is_dir}
                      name={file.name}
                      className={viewMode === "grid" ? "w-12 h-12" : "w-8 h-8"}
                    />
                  )}

                  <div className="min-w-0">
                    <p className="text-sm font-medium text-gray-200 truncate max-w-[120px] group-hover:text-blue-400 transition-colors">
                      {file.name}
                    </p>
                    <p className="text-xs text-gray-500 mt-1">
                      {file.is_dir
                        ? "Folder"
                        : formatDistanceToNow(new Date(file.mod_time), {
                            addSuffix: true,
                          })}
                    </p>
                  </div>
                </div>

                {/* Actions (Hover only) */}
                <div className="absolute top-2 right-2 opacity-0 group-hover:opacity-100 transition-opacity flex gap-1">
                  {!file.is_dir && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDownload(file.name);
                      }}
                      className="p-1.5 bg-gray-800 hover:bg-blue-600 rounded-md text-white shadow-lg"
                      title="Download"
                    >
                      <Download className="w-3.5 h-3.5" />
                    </button>
                  )}
                  <button
                    className="p-1.5 bg-gray-800 hover:bg-red-600 rounded-md text-white shadow-lg"
                    title="Delete"
                  >
                    <Trash2 className="w-3.5 h-3.5" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}