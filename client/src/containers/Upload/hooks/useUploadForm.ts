import { FileChunk, UploadRequest, UploadResponse } from "upload/upload_pb";
import { useState, useCallback, SyntheticEvent } from "react";
import { UploadClient } from "upload/UploadServiceClientPb";

export const useUploadForm = (client: UploadClient) => {
  const [file, setFile] = useState<File | null>(null);

  const onChange = useCallback(
    (event: SyntheticEvent<HTMLInputElement>) => {
      const target = event.target as HTMLInputElement;
      if (target.files && target.files[0]) {
        setFile(target.files[0]);
      }
    },
    [setFile]
  );

  const onSubmit = useCallback(
    async (event: SyntheticEvent) => {
      event.preventDefault();
      if (!file) return;

      const fileReader = new FileReader();
      const chunkSize = 1024 * 1024; // 1MB chunks (adjust as needed)
      let offset = 0;

      fileReader.onload = async function () {
        const chunkData = new Uint8Array(fileReader.result as ArrayBuffer);
        const totalChunks = Math.ceil(chunkData.length / chunkSize);

        for (let i = 0; i < totalChunks; i++) {
          const chunk = new FileChunk();
          chunk.setChunkNumber(i + 1);
          chunk.setData(chunkData.slice(offset, offset + chunkSize));
          offset += chunkSize;

          const uploadRequest = new UploadRequest();
          uploadRequest.setFileName(file.name);
          uploadRequest.setFileSize(file.size);
          uploadRequest.addChunks(chunk);
          uploadRequest.setIsLastChunk(i === totalChunks - 1);

          const response = await new Promise<UploadResponse>((resolve, reject) => {
            client.uploadFile(uploadRequest, {}, (err, res) => {
              if (err) {
                reject(err);
              } else {
                resolve(res);
              }
            });
          });

          console.log(response); // Handle response as needed
        }
      };

      fileReader.readAsArrayBuffer(file);
    },
    [file, client]
  );

  return {
    onChange,
    onSubmit
  };
};