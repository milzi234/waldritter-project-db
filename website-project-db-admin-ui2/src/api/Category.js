import apiClient from './api_client';

class CategoryAPI {
    constructor() {
        this.listeners = {}
        this.tagAPIs = {}
    }

    async getAll() {
        const response = await apiClient.get('/api/v1/categories');
        return response.data;
    }

    async get(id) {
        const response = await apiClient.get(`/api/v1/categories/${id}`);
        return response.data;
    }

    async create(data) {
        const response = await apiClient.post('/api/v1/categories', data);
        this.emit('created', response.data);
        return response.data;
    }

    async update(id, data) {
        const response = await apiClient.put(`/api/v1/categories/${id}`, data);
        this.emit('updated', response.data);
        return response.data;
    }

    async delete(id) {
        const response = await apiClient.delete(`/api/v1/categories/${id}`);
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
        this.basePath = '/api/v1/categories/' + categoryID + '/tags';
        this.categoryAPI = categoryAPI;
    }

    async getAll() {
        const response = await apiClient.get(this.basePath);
        return response.data;
    }

    async get(id) {
        const response = await apiClient.get(`${this.basePath}/${id}`);
        return response.data;
    }

    async create(data) {
        const response = await apiClient.post(this.basePath, data);
        this.categoryAPI.emit('tag-created', response.data);
        return response.data;
    }

    async update(id, data) {
        const response = await apiClient.put(`${this.basePath}/${id}`, data);
        this.categoryAPI.emit('tag-updated', response.data);
        return response.data;
    }

    async delete(id) {
        const response = await apiClient.delete(`${this.basePath}/${id}`);
        this.categoryAPI.emit('tag-deleted', response.data);
        return response.data;
    }
}

const categoryAPI = new CategoryAPI();
export default function useCategoryAPI() { return categoryAPI };
