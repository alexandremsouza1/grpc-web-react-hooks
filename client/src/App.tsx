import React from "react";
import { gRPCClients } from "gRPCClients";
import { UploadContainer } from "containers/Upload";

export const App = () => {
  return (
    <>
      <UploadContainer clients={gRPCClients} />
    </>
  );
};
