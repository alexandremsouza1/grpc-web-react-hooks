import React from "react";
import { Upload } from "components/Upload";
import { GRPCClients } from "gRPCClients";
import { useUploadForm } from "./hooks/useUploadForm";

type Props = {
  clients: GRPCClients;
};

export const UploadContainer: React.FC<Props> = ({ clients } : Props) => {
  const uploadClient = clients.uploadClient;
  const uploadFormState = useUploadForm(uploadClient);
  return (
    <div>
      <Upload {...uploadFormState} />
    </div>
  );
};
