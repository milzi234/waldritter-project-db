import axios from '@bundled-es-modules/axios/axios.js';
import { useToken } from '@/composables/auth_token';

class ProjectAPI {
    constructor() {
        this.api = axios;
        this.baseURL = 'http://localhost:3000/api/v1/projects';
        this.listeners = {}
        this.eventAPIs = {}
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

    async getNextOccurrence(id) {
        const { token } = useToken();
        const currentDate = new Date();
        const fiveYearsFromNow = new Date();
        fiveYearsFromNow.setFullYear(currentDate.getFullYear() + 5);
        const response = await this.api.get(`${this.baseURL}/${id}/occurrences?start_date=${currentDate.toISOString()}&end_date=${fiveYearsFromNow.toISOString()}&limit=1`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        if (response.data.length > 0) {
            return response.data[0];
        }
        return null;
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

    async uploadImage(id, formData) {
        const { token } = useToken();
        const response = await this.api.post(`${this.baseURL}/${id}/upload_image`, formData, {
            headers: {
                'Content-Type': 'multipart/form-data',
                Authorization: `Bearer ${token.value}`
            }
        });
        this.emit('image-uploaded', response.data);
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

    async getProjectTags(id) {
        const { token } = useToken();
        const response = await this.api.get(`${this.baseURL}/${id}/tags`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        return response.data;
    }

    async setProjectTags(id, tags) {
        const { token } = useToken();
        const response = await this.api.post(`${this.baseURL}/${id}/tags`, {tag_ids: tags}, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
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
        this.api = axios;
        this.baseURL = 'http://localhost:3000/api/v1/projects/' + projectID + '/events';
        this.projectAPI = projectAPI;
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
        this.projectAPI.emit('event-created', response.data);
        return response.data;
    }

    async update(id, data) {
        const { token } = useToken();
        const response = await this.api.put(`${this.baseURL}/${id}`, data, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.projectAPI.emit('event-updated', response.data);
        return response.data;
    }

    async delete(id) {
        const { token } = useToken();
        const response = await this.api.delete(`${this.baseURL}/${id}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.projectAPI.emit('event-deleted', response.data);
        return response.data;
    }

    async getOccurrences(id, startDate, endDate) {
        const { token } = useToken();
        const response = await this.api.get(`${this.baseURL}/${id}/occurrences?start_date=${startDate.toISOString()}&end_date=${endDate.toISOString()}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        return response.data;
    }

    async createException(id, occurrenceID) {
        const { token } = useToken();
        const response = await this.api.delete(`${this.baseURL}/${id}/occurrences/${occurrenceID}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.projectAPI.emit('exception-created', response.data);
        return response.data;
    }

    async getExceptions(id, startDate, endDate) {
        const { token } = useToken();
        const response = await this.api.get(`${this.baseURL}/${id}/exceptions?start_date=${startDate.toISOString()}&end_date=${endDate.toISOString()}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        return response.data;
    }

    async deleteException(id, exceptionID) {
        const { token } = useToken();
        const response = await this.api.delete(`${this.baseURL}/${id}/exceptions/${exceptionID}`, {
            headers: { Authorization: `Bearer ${token.value}` }
        });
        this.projectAPI.emit('exception-deleted', response.data);
        return response.data;
    }
}

const projectAPI = new ProjectAPI();
export default function useProjectAPI() { return projectAPI };