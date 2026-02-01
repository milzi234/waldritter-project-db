import axios from 'axios';

class ExtractionAPI {
    constructor() {
        this.api = axios;
        this.baseURL = 'http://localhost:3000/api/v1/extract';
    }

    getAuthToken() {
        const savedUser = localStorage.getItem('google_user');
        if (savedUser) {
            try {
                const user = JSON.parse(savedUser);
                return user.token;
            } catch (e) {
                return null;
            }
        }
        return null;
    }

    async extract(url) {
        const token = this.getAuthToken();
        const response = await this.api.post(this.baseURL, { url }, {
            headers: { Authorization: `Bearer ${token}` },
            timeout: 120000 // 2 minute timeout for exploration
        });
        return response.data;
    }
}

const extractionAPI = new ExtractionAPI();
export default function useExtractionAPI() { return extractionAPI };
