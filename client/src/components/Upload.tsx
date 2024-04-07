import React, { ChangeEvent, FormEvent } from "react";

type UploadProps = {
  onChange: (event: ChangeEvent<HTMLInputElement>) => void;
  onSubmit: (event: FormEvent<HTMLFormElement>) => void;
};

export const Upload: React.FC<UploadProps> = ({ onChange, onSubmit } : UploadProps) => {
  return (
    <div>
      <h1>Upload</h1>
      <form onSubmit={onSubmit}>
        <input type="file" onChange={onChange} />
        <button type="submit">Upload</button>
      </form> 
    </div>
  );
};