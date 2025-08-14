import axios from '@bundled-es-modules/axios/axios.js';
import { useToken } from '@/composables/auth_token';

class CategoryAPI {
    constructor() {
        this.api = axios;
        this.baseURL = 'http://localhost:3000/api/v1/categories';
        this.listeners = {}
        this.tagAPIs = {}
    }

    async getAll() {
        const { token } = useToken();
        const response = await this.api.get(this.baseURL, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        return response.data;
    }

    async get(id) {
        const { token } = useToken();
        const response = await this.api.get(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        return response.data;
    }

    async create(data) {
        const { token } = useToken();
        const response = await this.api.post(this.baseURL, data, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.emit('created', response.data);
        return response.data;
    }

    async update(id, data) {
        const { token } = useToken();
        const response = await this.api.put(`${this.baseURL}/${id}`, data, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.emit('updated', response.data);
        return response.data;
    }

    async delete(id) {
        const { token } = useToken();
        const response = await this.api.delete(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.emit('deleted', response.data);
        return response.data;
    }

    tagAPIFor(categoryID) {
        return this.tagAPIs[categoryID] || (this.tagAPIs[categoryID] = new TagAPI(categoryID, this));
    }

    on(event, callback) {
        if (!this.listeners[event]) {
            this.listeners[event] = [];
        }
        this.listeners[event].push(callback);
    }

    emit(event, data) {
        if (this.listeners[event]) {
            this.listeners[event].forEach(callback => {
                callback(data);
            });
        }
    }
}

class TagAPI {
    constructor(categoryID, categoryAPI) {
        this.api = axios;
        this.baseURL = 'http://localhost:3000/api/v1/categories/' + categoryID + '/tags';
        this.categoryAPI = categoryAPI;
    }

    async getAll() {
        const { token } = useToken();
        const response = await this.api.get(this.baseURL, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        return response.data;
    }

    async get(id) {
        const { token } = useToken();
        const response = await this.api.get(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        return response.data;
    }

    async create(data) {
        const { token } = useToken();
        const response = await this.api.post(this.baseURL, data, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.categoryAPI.emit('tag-created', response.data);
        return response.data;
    }

    async update(id, data) {
        const { token } = useToken();
        const response = await this.api.put(`${this.baseURL}/${id}`, data, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.categoryAPI.emit('tag-updated', response.data);
        return response.data;
    }

    async delete(id) {
        const { token } = useToken();
        const response = await this.api.delete(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.categoryAPI.emit('tag-deleted', response.data);
        return response.data;
    }
}

const categoryAPI = new CategoryAPI();
export default function useCategoryAPI() { return categoryAPI };