import { computed, onMounted, useTemplateRef, watch } from "../../vue"

export const Toast = {
  props: {
    id: String,
    title: String,
    message: String,
    type: string,
    show: Boolean,
  },
  components: {},
  setup: (props) => {
    const border = computed(
      () => props.type,
      () => {
        switch (props.type) {
          case "success":
            return "border-success"
          case "danger":
            return "border-danger"
          case "warning":
            return "border-warning"
          default:
            return "border-primary"
        }
      },
    )

    var el = useTemplateRef("liveToast")
    var bsEl = null

    onMounted(() => {
      bsEl = bootstrap.Toast.getOrCreateInstance(el.value)
    })

    watch(
      () => props.show,
      () => {
        if (props.show) {
          bsEl.show()
        } else {
          bsEl.hide()
        }
      },
    )

    return {
      el,
      border,
    }
  },
  template: /* HTML */ `
    <div class="toast-container position-fixed bottom-0 end-0 p-3">
      <div
        ref="el"
        :class="border"
        class="toast"
        role="alert"
        aria-live="assertive"
        aria-atomic="true"
      >
        <div class="toast-header">
          <strong class="me-auto">{{ title }}</strong>
          <button
            type="button"
            class="btn-close"
            data-bs-dismiss="toast"
            aria-label="Close"
          ></button>
        </div>
        <div class="toast-body">{{ message }}</div>
      </div>
    </div>
  `,
}
