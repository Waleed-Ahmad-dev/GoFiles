import {
  Folder,
  FileText,
  Image as ImageIcon,
  Music,
  Video,
  Code,
  Archive,
  File,
} from "lucide-react";

interface FileIconProps {
  isDir: boolean;
  name: string;
  className?: string;
}

export default function FileIcon({
  isDir,
  name,
  className = "w-6 h-6",
}: FileIconProps) {
  // Folders inherit the accent color from the parent component or default to blue
  // We use a specific utility here to make folders look distinct
  if (isDir)
    return <Folder className={`${className} text-sky-500 fill-sky-500/20`} />;

  const ext = name.split(".").pop()?.toLowerCase();

  switch (ext) {
    case "jpg":
    case "jpeg":
    case "png":
    case "gif":
    case "webp":
    case "svg":
      return <ImageIcon className={`${className} text-violet-500`} />;
    case "mp4":
    case "mkv":
    case "mov":
    case "avi":
      return <Video className={`${className} text-rose-500`} />;
    case "mp3":
    case "wav":
    case "ogg":
      return <Music className={`${className} text-pink-500`} />;
    case "zip":
    case "rar":
    case "7z":
    case "tar":
    case "gz":
      return <Archive className={`${className} text-amber-500`} />;
    case "js":
    case "ts":
    case "tsx":
    case "jsx":
    case "go":
    case "py":
    case "html":
    case "css":
    case "json":
      return <Code className={`${className} text-emerald-500`} />;
    case "txt":
    case "md":
    case "pdf":
    case "doc":
    case "docx":
      return <FileText className={`${className} text-slate-500`} />;
    default:
      return <File className={`${className} text-gray-400`} />;
  }
}