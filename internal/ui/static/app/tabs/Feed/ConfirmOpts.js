import { ky, TvStatus } from "../../utils.js"
import { computed, onMounted, ref, storeToRefs, watch } from "../../vue.js"
import { notify } from "../../components/Notify/Notify.js"
import { useTracked } from "../../stores/tracked.js"
import { useFeedStore } from "./feedStore.js"

const show = ref(false)
var resolver = null

const state = ref({
  title: "",
  message: "",
})

export const openConfirm = (title, message) => {
  state.value.title = title
  state.value.message = message
  show.value = true
  console.log(title, message)
  return new Promise((resolve) => {
    resolver = resolve
  })
}

window.openConfirm = openConfirm

export const ConfirmOpts = {
  props: {},
  setup(props) {
    const confirm = () => {
      if (resolver) {
        resolver(true)
        show.value = false
        resolver = null
      }
    }

    const cancel = () => {
      if (resolver) {
        resolver(false)
        show.value = false
        resolver = null
      }
    }

    var bSelf = null

    onMounted(() => {
      if (!bSelf) {
        const el = document.getElementById("confirmOpts")
        bSelf = bootstrap.Offcanvas.getOrCreateInstance(el)
        el.addEventListener("hidden.bs.offcanvas", () => {
          cancel()
        })
      }

      // el.addEventListener('show.bs.offcanvas', () => {
      // })
    })

    watch(show, () => {
      if (bSelf) {
        if (show.value) {
          bSelf.show()
        } else {
          bSelf.hide()
        }
      }
    })

    return {
      confirm,
      cancel,
      state,
      show,
    }
  },

  template: /* HTML */ `
    <div class="offcanvas offcanvas-end" tabindex="-1" id="confirmOpts">
      <div class="offcanvas-header border-bottom">
        <h5 class="offcanvas-title"><span class="text-muted">{{ state.title }}</span></h5>
        <button
          type="button"
          class="btn-close"
          data-bs-dismiss="offcanvas"
          aria-label="Close"
        ></button>
      </div>
      <div class="offcanvas-body p-0">
        <p class="px-3 mt-2">{{ state.message }}</p>
        <div class="list-group list-group-flush">
          <button
            class="list-group-item list-group-item-action btn btn-primary px-4 py-3 d-flex align-items-center border-0 text-primary"
            @click="confirm"
          >
            <i class="bi bi-plus-circle me-2"></i> Confirm
          </button>
          <button
            class="list-group-item list-group-item-action btn btn-primary px-4 py-3 d-flex align-items-center border-0 text-danger"
            @click="cancel"
          >
            <i class="bi bi-dash-circle me-2"></i> Cancel
          </button>
        </div>
      </div>
    </div>
  `,
}
