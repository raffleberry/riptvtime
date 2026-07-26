import { defineStore, ref } from "../../vue.js";

export const useConfirm = defineStore('ConfirmStore', () => {
    const loading = ref(false)
    const text = ref('')
    const isOpen = ref(false)
    const close = () => {
        isOpen.value = false
    }

    const onConfirm = ref(null)
    const onCancel = ref(null)

    const openDialog = (t = "", confirm = null, cancel = null) => {
        text.value = t
        onConfirm.value = confirm  || null
        onCancel.value = cancel  || null
        isOpen.value = true
    }

    const confirm = () => {
        if (onConfirm.value) onConfirm.value()
        close()
    }

    const cancel = () => {
        if (onCancel.value) onCancel.value()
        close()
    }


    return {
        // data
        text,
        loading,
        isOpen,

        // actions
        onConfirm,
        onCancel,
        openDialog,
        close,
        confirm,
        cancel,

    }

})



