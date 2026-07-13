import { PAGE, theme } from "../utils.js";
import { onMounted, ref } from "../vue.js";


const Series = {
    props: {

    },
    components: {

    },
    setup: (props) => {
        onMounted(() => {
        });

        const r = useRoute()

        const seriesId = ref(0)

        watch(() => r.params.id, (id) => {
            seriesId.value = id
            if (id) {
                fetchDetails(n)
            } else {
            }
        }, { immediate: true })


        return {
        }
    },
    template: `
    <div>
        <h1>Series</h1>
    </div>
    `
}
export { Series };

