import { isMobile, PAGE, routes, theme } from "../../utils.js"
import { computed, ref, useRoute } from "../../vue.js"
import { Toast } from "./Toast.js"

/**
 * Notify message types
 * @readonly
 * @enum {string}
 */
export const MsgType = {
  Success: "success",
  Error: "error",
  Warning: "warning",
  Info: "info",
}

const notifications = ref([])

/**
 * @param {MsgType} type
 * @param {string} message
 */
export const notify = (type, title, message) => {
  const id = Date.now()
  notifications.value.unshift({
    id: id,
    title: title,
    message: message,
    type: type,
  })

  setTimeout(() => {
    notifications.value = notifications.value.filter((n) => n.id !== id)
  }, 5500)
}

window.notify = notify
window.ntype = MsgType

export const Notify = {
  components: {
    Toast,
  },

  setup: (props) => {
    return {
      notifications,
    }
  },

  template: /* HTML */ `
    <div class="toast-container position-fixed top-0 end-0 p-3">
      <Toast
        v-for="n in notifications"
        :key="n.id"
        :id="n.id"
        :title="n.title"
        :message="n.message"
        :type="n.type"
      ></Toast>
    </div>
  `,
}
