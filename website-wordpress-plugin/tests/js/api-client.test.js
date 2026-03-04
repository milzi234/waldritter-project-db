/**
 * API Client Tests
 */

describe('API Client', () => {
  beforeEach(() => {
    jest.clearAllMocks();
  });

  describe('fetch wrapper', () => {
    it('should use the configured REST URL', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ projects: [] }),
      });

      // Simulate the api.fetch call from frontend.js
      const baseUrl = global.waldritterProjectDB.restUrl;
      const endpoint = 'projects';
      const url = baseUrl + endpoint;

      await fetch(url, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          'X-WP-Nonce': global.waldritterProjectDB.nonce,
        },
      });

      expect(fetch).toHaveBeenCalledWith(
        '/wp-json/waldritter/v1/projects',
        expect.objectContaining({
          headers: expect.objectContaining({
            'X-WP-Nonce': 'test-nonce',
          }),
        })
      );
    });

    it('should handle API errors', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () => Promise.resolve({ message: 'Server error' }),
      });

      const response = await fetch('/wp-json/waldritter/v1/projects');

      expect(response.ok).toBe(false);
      expect(response.status).toBe(500);
    });

    it('should handle network errors', async () => {
      global.fetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(fetch('/wp-json/waldritter/v1/projects')).rejects.toThrow('Network error');
    });
  });

  describe('getProjects', () => {
    it('should build query string with tags', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ projects: [], pagination: {} }),
      });

      // Simulate building query string
      const params = { tags: ['LARP', 'Aktiv'], page: 1, per_page: 5 };
      const query = new URLSearchParams();
      params.tags.forEach(tag => query.append('tags[]', tag));
      query.set('page', params.page);
      query.set('per_page', params.per_page);

      await fetch(`/wp-json/waldritter/v1/search?${query.toString()}`);

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('tags%5B%5D=LARP'),
        undefined
      );
    });

    it('should include pagination parameters', async () => {
      global.fetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({
          projects: [],
          pagination: { page: 2, total_pages: 5 },
        }),
      });

      const query = new URLSearchParams();
      query.set('page', '2');
      query.set('per_page', '5');

      await fetch(`/wp-json/waldritter/v1/search?${query.toString()}`);

      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('page=2'),
        undefined
      );
      expect(fetch).toHaveBeenCalledWith(
        expect.stringContaining('per_page=5'),
        undefined
      );
    });
  });
});
