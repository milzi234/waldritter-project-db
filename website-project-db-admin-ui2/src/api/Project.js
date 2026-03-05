import apiClient from './api_client';

class ProjectAPI {
    constructor() {
        this.listeners = {}
        this.eventAPIs = {}
    }

    async getAll() {
        const response = await apiClient.get('/api/v1/projects');
        return response.data;
    }

    async get(id) {
        const response = await apiClient.get(`/api/v1/projects/${id}`);
        return response.data;
    }

    async getNextOccurrence(id) {
        const currentDate = new Date();
        const fiveYearsFromNow = new Date();
        fiveYearsFromNow.setFullYear(currentDate.getFullYear() + 5);
        const response = await apiClient.get(`/api/v1/projects/${id}/occurrences?start_date=${currentDate.toISOString()}&end_date=${fiveYearsFromNow.toISOString()}&limit=1`);
        if (response.data.length > 0) {
            return response.data[0];
        }
        return null;
    }

    async create(data) {
        const response = await apiClient.post('/api/v1/projects', data);
        this.emit('created', response.data);
        return response.data;
    }

    async update(id, data) {
        const response = await apiClient.put(`/api/v1/projects/${id}`, data);
        this.emit('updated', response.data);
        return response.data;
    }

    async uploadImage(id, formData) {
        const response = await apiClient.post(`/api/v1/projects/${id}/upload_image`, formData, {
            headers: { 'Content-Type': 'multipart/form-data' }
        });
        this.emit('image-uploaded', response.data);
        return response.data;
    }

    async delete(id) {
        const response = await apiClient.delete(`/api/v1/projects/${id}`);
        this.emit('deleted', response.data);
        return response.data;
    }

    async getProjectTags(id) {
        const response = await apiClient.get(`/api/v1/projects/${id}/tags`);
        return response.data;
    }

    async setProjectTags(id, tags) {
        const response = await apiClient.post(`/api/v1/projects/${id}/tags`, {tag_ids: tags});
        this.emit('tags-updated', response.data);
        return response.data;
    }

    eventAPIFor(projectID) {
        return this.eventAPIs[projectID] || (this.eventAPIs[projectID] = new EventAPI(projectID, this));
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

class EventAPI {
    constructor(projectID, projectAPI) {
        this.basePath = '/api/v1/projects/' + projectID + '/events';
        this.projectAPI = projectAPI;
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
        this.projectAPI.emit('event-created', response.data);
        return response.data;
    }

    async update(id, data) {
        const response = await apiClient.put(`${this.basePath}/${id}`, data);
        this.projectAPI.emit('event-updated', response.data);
        return response.data;
    }

    async delete(id) {
        const response = await apiClient.delete(`${this.basePath}/${id}`);
        this.projectAPI.emit('event-deleted', response.data);
        return response.data;
    }

    async getOccurrences(id, startDate, endDate) {
        const response = await apiClient.get(`${this.basePath}/${id}/occurrences?start_date=${startDate.toISOString()}&end_date=${endDate.toISOString()}`);
        return response.data;
    }

    async createException(id, occurrenceID) {
        const response = await apiClient.delete(`${this.basePath}/${id}/occurrences/${occurrenceID}`);
        this.projectAPI.emit('exception-created', response.data);
        return response.data;
    }

    async getExceptions(id, startDate, endDate) {
        const response = await apiClient.get(`${this.basePath}/${id}/exceptions?start_date=${startDate.toISOString()}&end_date=${endDate.toISOString()}`);
        return response.data;
    }

    async deleteException(id, exceptionID) {
        const response = await apiClient.delete(`${this.basePath}/${id}/exceptions/${exceptionID}`);
        this.projectAPI.emit('exception-deleted', response.data);
        return response.data;
    }
}

const projectAPI = new ProjectAPI();
export default function useProjectAPI() { return projectAPI };
