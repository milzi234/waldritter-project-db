import axios from 'axios';

class CategoryAPI {
    constructor() {
        this.api = axios;
        this.baseURL = 'http://localhost:3000/api/v1/categories';
        this.listeners = {}
        this.tagAPIs = {}
    }

    getAuthToken() {
        // Get token from localStorage directly
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

    async getAll() {
        const token = this.getAuthToken();
        const response = await this.api.get(this.baseURL, {
            headers: { Authorization: `Bearer ${token}` }
        });
        return response.data;
    }

    async get(id) {
        const token = this.getAuthToken();
        const response = await this.api.get(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        return response.data;
    }

    async create(data) {
        const token = this.getAuthToken();
        const response = await this.api.post(this.baseURL, data, {
            headers: { Authorization: `Bearer ${token}` }
        });
        this.emit('created', response.data);
        return response.data;
    }

    async update(id, data) {
        const token = this.getAuthToken();
        const response = await this.api.put(`${this.baseURL}/${id}`, data, {
            headers: { Authorization: `Bearer ${token}` }
        });
        this.emit('updated', response.data);
        return response.data;
    }

    async delete(id) {
        const token = this.getAuthToken();
        const response = await this.api.delete(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token}` }
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

    getAuthToken() {
        // Get token from localStorage directly
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

    async getAll() {
        const token = this.getAuthToken();
        const response = await this.api.get(this.baseURL, {
            headers: { Authorization: `Bearer ${token}` }
        });
        return response.data;
    }

    async get(id) {
        const token = this.getAuthToken();
        const response = await this.api.get(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        return response.data;
    }

    async create(data) {
        const token = this.getAuthToken();
        const response = await this.api.post(this.baseURL, data, {
            headers: { Authorization: `Bearer ${token}` }
        });
        this.categoryAPI.emit('tag-created', response.data);
        return response.data;
    }

    async update(id, data) {
        const token = this.getAuthToken();
        const response = await this.api.put(`${this.baseURL}/${id}`, data, {
            headers: { Authorization: `Bearer ${token}` }
        });
        this.categoryAPI.emit('tag-updated', response.data);
        return response.data;
    }

    async delete(id) {
        const token = this.getAuthToken();
        const response = await this.api.delete(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token}` }
        });
        this.categoryAPI.emit('tag-deleted', response.data);
        return response.data;
    }
}

const categoryAPI = new CategoryAPI();
export default function useCategoryAPI() { return categoryAPI };