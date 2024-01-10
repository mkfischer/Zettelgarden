import React, { useState, useEffect } from "react";
import {
  renderFile,
  uploadFile,
  getAllFiles,
  deleteFile,
  editFile,
} from "../api";
import { sortCards } from "../utils";
import { FileRenameModal } from "./FileRenameModal.js";
import { FileListItem } from "./FileListItem";

export function FileVault({ handleViewCard }) {
  const [files, setFiles] = useState([]);
  const [isRenameModalOpen, setIsRenameModalOpen] = useState(false);
  const [fileToRename, setFileToRename] = useState(null);

  const openRenameModal = (file) => {
    setFileToRename(file);
    setIsRenameModalOpen(true);
  };

  function onDelete(file_id) {
    setFiles(files.filter((file) => file.id !== file_id));
  }

  function onRename(fileId, updatedFile) {
    setFiles((prevFiles) =>
      prevFiles.map((f) => (f.id === updatedFile.id ? updatedFile : f)),
    );
    setIsRenameModalOpen(false);
  }

  useEffect(() => {
    getAllFiles().then((data) => setFiles(sortCards(data, "sortNewOld")));
  }, []);
  return (
    <>
      <FileRenameModal
        isOpen={isRenameModalOpen}
        onClose={() => setIsRenameModalOpen(false)}
        onRename={onRename}
        file={fileToRename}
      />
      <h3>File Vault</h3>
      <ul>
        {files &&
          files.map((file, index) => (
            <FileListItem
              file={file}
              onDelete={onDelete}
              handleViewCard={handleViewCard}
              openRenameModal={openRenameModal}
              displayCard={true}
            />
          ))}
      </ul>
    </>
  );
}