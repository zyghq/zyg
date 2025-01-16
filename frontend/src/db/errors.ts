export class HTTPError extends Error {
  status: number;
  statusText: string;

  constructor({
    message,
    status,
    statusText,
  }: {
    message: string;
    status: number;
    statusText: string;
  }) {
    super(message);
    this.status = status;
    this.statusText = statusText;
    this.name = "HTTPError";
  }
}
