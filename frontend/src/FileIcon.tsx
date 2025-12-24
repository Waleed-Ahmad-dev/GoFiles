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
  if (isDir)
    return <Folder className={`${className} text-blue-400 fill-blue-400/20`} />;

  const ext = name.split(".").pop()?.toLowerCase();

  switch (ext) {
    case "jpg":
    case "jpeg":
    case "png":
    case "gif":
    case "webp":
      return <ImageIcon className={`${className} text-purple-400`} />;
    case "mp4":
    case "mkv":
    case "mov":
      return <Video className={`${className} text-red-400`} />;
    case "mp3":
    case "wav":
      return <Music className={`${className} text-pink-400`} />;
    case "zip":
    case "rar":
    case "7z":
    case "tar":
    case "gz":
      return <Archive className={`${className} text-yellow-400`} />;
    case "js":
    case "ts":
    case "tsx":
    case "jsx":
    case "go":
    case "py":
    case "html":
    case "css":
    case "json":
      return <Code className={`${className} text-green-400`} />;
    case "txt":
    case "md":
      return <FileText className={`${className} text-gray-400`} />;
    default:
      return <File className={`${className} text-gray-500`} />;
  }
}