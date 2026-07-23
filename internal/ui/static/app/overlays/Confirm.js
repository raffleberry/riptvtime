import { useConfirm } from "../stores/confirm.js"
import { TvStatus } from "../utils.js"
import { computed, onMounted, ref, storeToRefs, watch } from "../vue.js"

const Confirm = {

    setup() {

        const store = useConfirm()

        const { loading, text, isOpen } = storeToRefs(store)
        const { confirm, cancel } = store

        var el = null
        

        watch(isOpen, () => {
            if (isOpen.value) el.show()
            else el.hide()
        })
      
        onMounted(() => {
            let elDom = document.getElementById('confirmDialog')
            elDom.addEventListener('hide.bs.modal', () => {
                cancel()
            })

            el = new bootstrap.Modal(elDom, {})

        })

        return {
            text,
            confirm,
            cancel,
        }

    },

    template: `
    <div class="modal fade" id="confirmDialog" tabindex="-1" >
        <div class="modal-dialog">
            <div class="modal-content">
                <div class="modal-header">
                    <h5 class="modal-title">Confirmation</h5>
                    <button type="button" class="btn-close" aria-label="Close" @click="cancel"></button>
                </div>
                <div class="modal-body">
                    <p>{{ text }}</p>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-secondary" @click="cancel">Close</button>
                    <button type="button" class="btn btn-primary" @click="confirm">Confirm</button>
                </div>
            </div>
        </div>
    </div>
`
}

export { Confirm }