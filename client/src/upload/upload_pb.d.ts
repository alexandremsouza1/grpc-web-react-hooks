import * as jspb from 'google-protobuf'



export class FileChunk extends jspb.Message {
  getChunkNumber(): number;
  setChunkNumber(value: number): FileChunk;

  getData(): Uint8Array | string;
  getData_asU8(): Uint8Array;
  getData_asB64(): string;
  setData(value: Uint8Array | string): FileChunk;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): FileChunk.AsObject;
  static toObject(includeInstance: boolean, msg: FileChunk): FileChunk.AsObject;
  static serializeBinaryToWriter(message: FileChunk, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): FileChunk;
  static deserializeBinaryFromReader(message: FileChunk, reader: jspb.BinaryReader): FileChunk;
}

export namespace FileChunk {
  export type AsObject = {
    chunkNumber: number,
    data: Uint8Array | string,
  }
}

export class UploadRequest extends jspb.Message {
  getFileName(): string;
  setFileName(value: string): UploadRequest;

  getFileSize(): number;
  setFileSize(value: number): UploadRequest;

  getChunksList(): Array<FileChunk>;
  setChunksList(value: Array<FileChunk>): UploadRequest;
  clearChunksList(): UploadRequest;
  addChunks(value?: FileChunk, index?: number): FileChunk;

  getIsLastChunk(): boolean;
  setIsLastChunk(value: boolean): UploadRequest;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadRequest.AsObject;
  static toObject(includeInstance: boolean, msg: UploadRequest): UploadRequest.AsObject;
  static serializeBinaryToWriter(message: UploadRequest, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadRequest;
  static deserializeBinaryFromReader(message: UploadRequest, reader: jspb.BinaryReader): UploadRequest;
}

export namespace UploadRequest {
  export type AsObject = {
    fileName: string,
    fileSize: number,
    chunksList: Array<FileChunk.AsObject>,
    isLastChunk: boolean,
  }
}

export class UploadResponse extends jspb.Message {
  getMessage(): string;
  setMessage(value: string): UploadResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UploadResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UploadResponse): UploadResponse.AsObject;
  static serializeBinaryToWriter(message: UploadResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UploadResponse;
  static deserializeBinaryFromReader(message: UploadResponse, reader: jspb.BinaryReader): UploadResponse;
}

export namespace UploadResponse {
  export type AsObject = {
    message: string,
  }
}

