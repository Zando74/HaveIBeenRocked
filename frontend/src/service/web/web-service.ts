class WebService {
    private baseUrl: string;

    constructor(
        baseUrl: string,
      ) {
        this.baseUrl = baseUrl;
      }

      async post(endpoint: string, data: unknown): Promise<Response> {
        const response = await fetch(`${this.baseUrl}${endpoint}`, {
          method: "POST",
          headers: {"Content-Type": "application/json"},
          body: JSON.stringify(data),
        });
    
        return response;
      }
}

export default WebService;