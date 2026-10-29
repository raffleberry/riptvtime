const Empty = {
  components: {},
  props: {
    icon: {
      type: String,
      default: "bi bi-calendar-check-fill",
    },
    title: {
      type: String,
      default: "No Data Available",
    },
    message: {
      type: String,
      default:
        "There is currently no information to display here. Try adding a new record or refreshing the page.",
    },
  },
  setup() {
    return {}
  },
  template: /* HTML */ `
    <div class="d-flex flex-grow-1 align-items-center justify-content-center h-100 my-auto p-3">
      <div
        class="card border-0 shadow-sm text-center d-flex flex-column justify-content-center p-5 w-100 h-100"
        style="max-width: 1200px; max-height: 80vh; border-radius: 16px;"
      >
        <div class="card-body d-flex flex-column justify-content-center align-items-center py-5">
          <div class="mb-4 opacity-50">
            <i :class="icon" style="font-size: 8rem;"></i>
          </div>

          <h2 class="card-title fw-bold mb-3 fs-3 fs-md-2">{{ title }}</h2>
          <p class="card-text text-muted fs-6 mb-0" style="max-width: 550px;">{{ message }}</p>
        </div>
      </div>
    </div>
  `,
}

export { Empty }
