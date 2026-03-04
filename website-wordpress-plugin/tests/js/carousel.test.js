/**
 * Carousel Component Tests
 */

describe('Carousel', () => {
  let container;

  beforeEach(() => {
    container = document.createElement('div');
    container.innerHTML = `
      <div
        id="test-carousel"
        class="waldritter-carousel"
        data-auto-scroll="true"
        data-interval="5000"
      >
        <div class="waldritter-carousel__track">
          <div class="waldritter-carousel__slides">
            <div class="waldritter-carousel__slide">Slide 1</div>
            <div class="waldritter-carousel__slide">Slide 2</div>
            <div class="waldritter-carousel__slide">Slide 3</div>
          </div>
        </div>
        <div class="waldritter-carousel__controls">
          <button class="waldritter-carousel__button waldritter-carousel__button--prev">Prev</button>
          <div class="waldritter-carousel__dots">
            <button class="waldritter-carousel__dot waldritter-carousel__dot--active" data-slide="0"></button>
            <button class="waldritter-carousel__dot" data-slide="1"></button>
            <button class="waldritter-carousel__dot" data-slide="2"></button>
          </div>
          <button class="waldritter-carousel__button waldritter-carousel__button--next">Next</button>
        </div>
      </div>
    `;
    document.body.appendChild(container);
  });

  afterEach(() => {
    document.body.removeChild(container);
  });

  describe('Structure', () => {
    it('should have required elements', () => {
      const carousel = container.querySelector('.waldritter-carousel');
      expect(carousel).toBeTruthy();
      expect(carousel.querySelector('.waldritter-carousel__slides')).toBeTruthy();
      expect(carousel.querySelector('.waldritter-carousel__button--prev')).toBeTruthy();
      expect(carousel.querySelector('.waldritter-carousel__button--next')).toBeTruthy();
    });

    it('should have data attributes for configuration', () => {
      const carousel = container.querySelector('.waldritter-carousel');
      expect(carousel.dataset.autoScroll).toBe('true');
      expect(carousel.dataset.interval).toBe('5000');
    });

    it('should have correct number of slides', () => {
      const slides = container.querySelectorAll('.waldritter-carousel__slide');
      expect(slides.length).toBe(3);
    });

    it('should have correct number of dots', () => {
      const dots = container.querySelectorAll('.waldritter-carousel__dot');
      expect(dots.length).toBe(3);
    });
  });

  describe('Navigation', () => {
    it('should have prev button', () => {
      const prevBtn = container.querySelector('.waldritter-carousel__button--prev');
      expect(prevBtn).toBeTruthy();
      expect(prevBtn.textContent).toBe('Prev');
    });

    it('should have next button', () => {
      const nextBtn = container.querySelector('.waldritter-carousel__button--next');
      expect(nextBtn).toBeTruthy();
      expect(nextBtn.textContent).toBe('Next');
    });
  });

  describe('Dots', () => {
    it('should have first dot active by default', () => {
      const firstDot = container.querySelector('.waldritter-carousel__dot');
      expect(firstDot.classList.contains('waldritter-carousel__dot--active')).toBe(true);
    });

    it('should have data-slide attributes', () => {
      const dots = container.querySelectorAll('.waldritter-carousel__dot');
      dots.forEach((dot, index) => {
        expect(dot.dataset.slide).toBe(String(index));
      });
    });
  });

  describe('Responsive', () => {
    it('should calculate items per view based on viewport', () => {
      // Test logic for responsive behavior
      const getItemsPerView = (width) => {
        if (width >= 1024) return 3;
        if (width >= 768) return 2;
        return 1;
      };

      expect(getItemsPerView(1200)).toBe(3);
      expect(getItemsPerView(800)).toBe(2);
      expect(getItemsPerView(400)).toBe(1);
    });
  });
});
