import apiClient from './api_client';

class ExtractionAPI {
    async extract(url) {
        const response = await apiClient.post('/api/v1/extract', { url }, {
            timeout: 120000 // 2 minute timeout for exploration
        });
        return response.data;
    }

    async extractFromText(text) {
        const response = await apiClient.post('/api/v1/extract_text', { text }, {
            timeout: 120000
        });
        return response.data;
    }

    async generateImage(title, description, keywords = [], variation = 0) {
        const response = await apiClient.post(
            '/api/v1/generate_image',
            { title, description, keywords, variation },
            { timeout: 60000 }
        );
        return response.data;
    }
}

const extractionAPI = new ExtractionAPI();
export default function useExtractionAPI() { return extractionAPI };
