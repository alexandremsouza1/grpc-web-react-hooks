import { UploadClient } from "upload/UploadServiceClientPb";
import { MessengerClient } from "messenger/MessengerServiceClientPb";

export type GRPCClients = {
  uploadClient: UploadClient;
  messengerClient: MessengerClient;
};

export const gRPCClients = {
  uploadClient: new UploadClient(`http://localhost:8080`),
  messengerClient: new MessengerClient(`http://localhost:8080`)
};
