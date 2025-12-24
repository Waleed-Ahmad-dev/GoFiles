/* eslint-disable react-hooks/exhaustive-deps */
import { useEffect, useState, useRef } from "react";
import { api, type FileInfo } from "../Utils/api";
import FileIcon from "./FileIcon";
import { useTheme, type AccentColor } from "../Context/ThemeContext";
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
  Settings,
  ChevronRight,
} from "lucide-react";
import { formatDistanceToNow } from "date-fns";

export default function Dashboard() {
  const [files, setFiles] = useState<FileInfo[]>([]);
  const [currentPath, setCurrentPath] = useState("");
  const [loading, setLoading] = useState(false);
  const [viewMode, setViewMode] = useState<"grid" | "list">("grid");
  const [searchQuery, setSearchQuery] = useState("");
  const [isSettingsOpen, setIsSettingsOpen] = useState(false);

  // Theme Hooks
  const { theme, setTheme, accent, setAccent, accentStyles } = useTheme();
  const settingsRef = useRef<HTMLDivElement>(null);

  // Close settings on click outside
  useEffect(() => {
    function handleClickOutside(event: MouseEvent) {
      if (
        settingsRef.current &&
        !settingsRef.current.contains(event.target as Node)
      ) {
        setIsSettingsOpen(false);
      }
    }
    document.addEventListener("mousedown", handleClickOutside);
    return () => document.removeEventListener("mousedown", handleClickOutside);
  }, []);

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
    const downloadUrl = `${api.defaults.baseURL}/download?path=${
      currentPath ? currentPath + "/" : ""
    }${fileName}`;
    window.open(downloadUrl, "_blank");
  };

  // Filter files based on search
  const filteredFiles = files.filter((f) =>
    f.name.toLowerCase().includes(searchQuery.toLowerCase())
  );

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-950 flex flex-col transition-colors duration-300 font-sans">
      {/* --- HEADER --- */}
      <header className="sticky top-0 z-20 backdrop-blur-md bg-white/80 dark:bg-gray-900/80 border-b border-gray-200 dark:border-gray-800 transition-all duration-300">
        <div className="px-6 py-3 flex items-center justify-between">
          {/* Logo & Breadcrumbs */}
          <div className="flex items-center gap-6 overflow-hidden">
            <div className="flex items-center gap-2 font-bold text-xl tracking-tight text-gray-900 dark:text-white shrink-0">
              <span
                className={`bg-clip-text text-transparent bg-linear-to-r from-gray-900 to-gray-600 dark:from-white dark:to-gray-400`}
              >
                GoFiles
              </span>
            </div>

            <div className="hidden md:flex items-center text-sm text-gray-500 dark:text-gray-400 overflow-hidden whitespace-nowrap mask-image-linear-to-r">
              <button
                onClick={() => setCurrentPath("")}
                className="hover:text-gray-900 dark:hover:text-white transition-colors p-1 rounded"
              >
                <Home className="w-4 h-4" />
              </button>
              {currentPath && (
                <>
                  <ChevronRight className="w-4 h-4 mx-1 opacity-50" />
                  <span className={`font-medium ${accentStyles.text} truncate`}>
                    {currentPath.replace(/\//g, " / ")}
                  </span>
                </>
              )}
            </div>
          </div>

          {/* Right Controls */}
          <div className="flex items-center gap-3">
            {/* Search */}
            <div className="relative hidden md:block group">
              <Search className="w-4 h-4 absolute left-3 top-1/2 -translate-y-1/2 text-gray-400 group-focus-within:text-gray-600 dark:group-focus-within:text-gray-200 transition-colors" />
              <input
                type="text"
                placeholder="Search files..."
                value={searchQuery}
                onChange={(e) => setSearchQuery(e.target.value)}
                className={`bg-gray-100 dark:bg-gray-800 border-transparent focus:bg-white dark:focus:bg-gray-950 border focus:border-gray-200 dark:focus:border-gray-700 text-sm rounded-full pl-9 pr-4 py-1.5 w-48 focus:w-64 transition-all outline-none ${accentStyles.ring}`}
              />
            </div>

            <div className="h-6 w-px bg-gray-200 dark:bg-gray-800 mx-1 hidden sm:block"></div>

            {/* View Toggle */}
            <div className="bg-gray-100 dark:bg-gray-800 rounded-lg p-1 flex">
              <button
                onClick={() => setViewMode("grid")}
                className={`p-1.5 rounded-md transition-all ${
                  viewMode === "grid"
                    ? "bg-white dark:bg-gray-700 shadow-sm text-gray-900 dark:text-white"
                    : "text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
                }`}
              >
                <Grid className="w-4 h-4" />
              </button>
              <button
                onClick={() => setViewMode("list")}
                className={`p-1.5 rounded-md transition-all ${
                  viewMode === "list"
                    ? "bg-white dark:bg-gray-700 shadow-sm text-gray-900 dark:text-white"
                    : "text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
                }`}
              >
                <ListIcon className="w-4 h-4" />
              </button>
            </div>

            {/* Settings Dropdown */}
            <div className="relative" ref={settingsRef}>
              <button
                onClick={() => setIsSettingsOpen(!isSettingsOpen)}
                className={`p-2 text-gray-500 hover:bg-gray-100 dark:text-gray-400 dark:hover:bg-gray-800 rounded-full transition-colors ${
                  isSettingsOpen
                    ? "bg-gray-100 dark:bg-gray-800 text-gray-900 dark:text-white"
                    : ""
                }`}
              >
                <Settings className="w-5 h-5" />
              </button>

              {isSettingsOpen && (
                <div className="absolute right-0 mt-3 w-64 bg-white dark:bg-gray-900 rounded-xl shadow-2xl border border-gray-200 dark:border-gray-800 p-4 transform origin-top-right animate-in fade-in zoom-in-95 duration-200">
                  <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">
                    Appearance
                  </h3>

                  {/* Theme Select */}
                  <div className="bg-gray-100 dark:bg-gray-800 rounded-lg p-1 flex mb-4">
                    {(["light", "system", "dark"] as const).map((t) => (
                      <button
                        key={t}
                        onClick={() => setTheme(t)}
                        className={`flex-1 py-1.5 text-xs font-medium capitalize rounded-md transition-all ${
                          theme === t
                            ? "bg-white dark:bg-gray-700 shadow-sm text-gray-900 dark:text-white"
                            : "text-gray-500 hover:text-gray-700 dark:hover:text-gray-300"
                        }`}
                      >
                        {t}
                      </button>
                    ))}
                  </div>

                  <h3 className="text-xs font-semibold text-gray-400 uppercase tracking-wider mb-3">
                    Accent Color
                  </h3>
                  <div className="grid grid-cols-6 gap-2 mb-4">
                    {(
                      [
                        "blue",
                        "violet",
                        "emerald",
                        "rose",
                        "amber",
                        "cyan",
                      ] as AccentColor[]
                    ).map((c) => (
                      <button
                        key={c}
                        onClick={() => setAccent(c)}
                        className={`w-6 h-6 rounded-full flex items-center justify-center transition-transform hover:scale-110 ${
                          c === accent
                            ? "ring-2 ring-offset-2 ring-offset-white dark:ring-offset-gray-900 ring-gray-400"
                            : ""
                        }`}
                        style={{
                          backgroundColor: `var(--color-${c}-500, ${
                            c === "blue"
                              ? "#3b82f6"
                              : c === "violet"
                              ? "#8b5cf6"
                              : c === "emerald"
                              ? "#10b981"
                              : c === "rose"
                              ? "#f43f5e"
                              : c === "amber"
                              ? "#f59e0b"
                              : "#06b6d4"
                          })`,
                        }}
                      >
                        {/* Fallback inline styles for preview dots if tailwind classes aren't enough */}
                        <div
                          className={`w-full h-full rounded-full ${
                            c === "blue"
                              ? "bg-blue-500"
                              : c === "violet"
                              ? "bg-violet-500"
                              : c === "emerald"
                              ? "bg-emerald-500"
                              : c === "rose"
                              ? "bg-rose-500"
                              : c === "amber"
                              ? "bg-amber-500"
                              : "bg-cyan-500"
                          }`}
                        ></div>
                      </button>
                    ))}
                  </div>

                  <div className="border-t border-gray-200 dark:border-gray-800 pt-3">
                    <button
                      onClick={() => window.location.reload()}
                      className="w-full flex items-center justify-center gap-2 px-3 py-2 text-sm text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20 rounded-lg transition-colors"
                    >
                      <LogOut className="w-4 h-4" /> Sign Out
                    </button>
                  </div>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* --- TOOLBAR --- */}
        <div className="px-6 py-2 border-t border-gray-200 dark:border-gray-800 bg-gray-50/50 dark:bg-gray-900/50 flex items-center gap-2 overflow-x-auto">
          <button
            onClick={handleGoUp}
            disabled={!currentPath}
            className="toolbar-btn"
          >
            <ArrowUp className="w-4 h-4" /> Up
          </button>
          <button onClick={loadFiles} className="toolbar-btn">
            <RefreshCw className={`w-4 h-4 ${loading ? "animate-spin" : ""}`} />{" "}
            Refresh
          </button>
          <div className="w-px h-4 bg-gray-300 dark:bg-gray-700 mx-2"></div>
          <button
            className={`toolbar-btn ${accentStyles.text} hover:${accentStyles.lightBg}`}
          >
            <FolderPlus className="w-4 h-4" /> New Folder
          </button>
          <button
            className={`toolbar-btn ${accentStyles.text} hover:${accentStyles.lightBg}`}
          >
            <FilePlus className="w-4 h-4" /> Upload
          </button>
        </div>
      </header>

      {/* --- MAIN CONTENT --- */}
      <main className="flex-1 p-6 overflow-y-auto scroll-smooth">
        {filteredFiles.length === 0 && !loading ? (
          <div className="flex flex-col items-center justify-center h-[50vh] text-gray-400">
            <div
              className={`p-6 rounded-full bg-gray-100 dark:bg-gray-900 mb-4`}
            >
              <Folder className="w-12 h-12 text-gray-300 dark:text-gray-700" />
            </div>
            <p className="font-medium">This folder is empty</p>
          </div>
        ) : (
          <div
            className={
              viewMode === "grid"
                ? "grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4"
                : "flex flex-col gap-1"
            }
          >
            {filteredFiles.map((file) => (
              <div
                key={file.name}
                onClick={() => file.is_dir && handleNavigate(file.name)}
                className={`
                  group relative border rounded-xl transition-all cursor-pointer select-none
                  ${
                    viewMode === "grid"
                      ? "aspect-4/5 p-4 flex flex-col items-center text-center justify-between hover:-translate-y-1"
                      : "p-3 flex items-center justify-between hover:translate-x-1"
                  }
                  border-gray-200 dark:border-gray-800 bg-white dark:bg-gray-900/50
                  hover:shadow-lg dark:hover:shadow-black/40 hover:border-gray-300 dark:hover:border-gray-700
                `}
              >
                {/* Visual Selection Indicator on Hover */}
                <div
                  className={`absolute inset-0 rounded-xl border-2 opacity-0 group-hover:opacity-100 transition-opacity pointer-events-none ${accentStyles.border}`}
                ></div>

                <div
                  className={`flex items-center w-full ${
                    viewMode === "grid"
                      ? "flex-col gap-4 flex-1 justify-center"
                      : "gap-4"
                  }`}
                >
                  {/* Thumbnail / Icon */}
                  <div className="relative">
                    {["jpg", "jpeg", "png", "gif", "webp"].includes(
                      file.type.replace(".", "").toLowerCase()
                    ) ? (
                      <img
                        src={`${api.defaults.baseURL}/thumbnail?path=${
                          currentPath ? currentPath + "/" : ""
                        }${file.name}`}
                        alt={file.name}
                        className={`object-cover rounded-lg shadow-sm ${
                          viewMode === "grid" ? "w-20 h-20" : "w-10 h-10"
                        }`}
                        loading="lazy"
                      />
                    ) : (
                      <FileIcon
                        isDir={file.is_dir}
                        name={file.name}
                        className={
                          viewMode === "grid" ? "w-16 h-16" : "w-10 h-10"
                        }
                      />
                    )}
                  </div>

                  <div className="min-w-0 flex-1">
                    <p
                      className={`font-medium text-gray-700 dark:text-gray-200 truncate ${
                        viewMode === "grid"
                          ? "text-sm w-full px-2"
                          : "text-base"
                      }`}
                    >
                      {file.name}
                    </p>
                    <p className="text-xs text-gray-400 mt-1">
                      {file.is_dir
                        ? "Folder"
                        : formatDistanceToNow(new Date(file.mod_time), {
                            addSuffix: true,
                          })}
                    </p>
                  </div>
                </div>

                {/* Actions */}
                <div
                  className={`flex gap-1 ${
                    viewMode === "grid"
                      ? "opacity-0 group-hover:opacity-100 absolute top-2 right-2"
                      : "opacity-0 group-hover:opacity-100 transition-opacity"
                  }`}
                >
                  {!file.is_dir && (
                    <button
                      onClick={(e) => {
                        e.stopPropagation();
                        handleDownload(file.name);
                      }}
                      className="action-btn"
                    >
                      <Download className="w-4 h-4" />
                    </button>
                  )}
                  <button className="action-btn text-red-500 hover:bg-red-50 dark:hover:bg-red-900/20">
                    <Trash2 className="w-4 h-4" />
                  </button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>

      {/* Footer / Status Bar */}
      <footer className="bg-white dark:bg-gray-900 border-t border-gray-200 dark:border-gray-800 px-6 py-2 text-xs text-gray-500 flex justify-between">
        <span>
          {filteredFiles.length} item{filteredFiles.length !== 1 && "s"}
        </span>
        <span>{loading ? "Syncing..." : "Ready"}</span>
      </footer>

      {/* Inline Styles for utility usage simplification */}
      <style>{`
        .toolbar-btn {
          @apply flex items-center gap-2 px-3 py-1.5 text-xs font-medium text-gray-600 dark:text-gray-400 hover:bg-gray-200 dark:hover:bg-gray-800 rounded-md transition-colors disabled:opacity-50;
        }
        .action-btn {
          @apply p-2 bg-white dark:bg-gray-800 shadow-sm border border-gray-200 dark:border-gray-700 rounded-lg text-gray-600 dark:text-gray-300 hover:text-gray-900 dark:hover:text-white transition-colors;
        }
      `}</style>
    </div>
  );
}
