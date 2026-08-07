import { computed, onMounted, useTemplateRef, watch } from "../../vue.js"

export const Toast = {
  props: {
    id: String,
    title: String,
    message: String,
    type: String,
  },
  components: {},
  setup: (props) => {
    const border = () => {
      switch (props.type) {
        case "success":
          return "border border-success"
        case "error":
          return "border border-danger"
        case "warning":
          return "border border-warning"
        default:
          return "border border-primary"
      }
    }

    var el = useTemplateRef("liveToast")
    var bsEl = null

    onMounted(() => {
      bsEl = bootstrap.Toast.getOrCreateInstance(el.value)
      bsEl.show()
    })

    return {
      el,
      border,
    }
  },
  template: /* HTML */ `
    <div
      ref="el"
      :class="border()"
      class="toast"
      role="alert"
      aria-live="assertive"
      aria-atomic="true"
    >
      <div class="toast-header">
        <strong class="me-auto">{{ title }}</strong>
        <button type="button" class="btn-close" data-bs-dismiss="toast" aria-label="Close"></button>
      </div>
      <div class="toast-body">{{ message }}</div>
    </div>
  `,
}
