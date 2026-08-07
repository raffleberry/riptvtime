import { isMobile, PAGE, theme, routes } from "../../utils.js"
import { computed, ref, useRoute } from "../../vue.js"

export const notify = (message) => {
  console.error("Error: ", message)
}

window.notify = notify
