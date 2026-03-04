import axios from 'axios';

class ExtractionAPI {
    constructor() {
        this.api = axios;
        this.baseURL = (import.meta.env.VITE_API_BASE_URL || 'http://localhost:3000') + '/api/v1';
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
        const response = await this.api.post(`${this.baseURL}/extract`, { url }, {
            headers: { Authorization: `Bearer ${token}` },
            timeout: 120000 // 2 minute timeout for exploration
        });
        return response.data;
    }

    async extractFromText(text) {
        const token = this.getAuthToken();
        const response = await this.api.post(`${this.baseURL}/extract_text`, { text }, {
            headers: { Authorization: `Bearer ${token}` },
            timeout: 120000
        });
        return response.data;
    }

    async generateImage(title, description, keywords = [], variation = 0) {
        const token = this.getAuthToken();
        const response = await this.api.post(
            `${this.baseURL}/generate_image`,
            { title, description, keywords, variation },
            {
                headers: { Authorization: `Bearer ${token}` },
                timeout: 60000 // 1 minute timeout for image generation
            }
        );
        return response.data;
    }
}

const extractionAPI = new ExtractionAPI();
export default function useExtractionAPI() { return extractionAPI };
