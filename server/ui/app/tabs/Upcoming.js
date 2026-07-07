import { currentPage, PAGE, theme, updatePage } from "../utils.js";
import { onMounted, ref } from "../vue.js";


const Upcoming = {
    props: {

    },
    components: {

    },
    setup: (props) => {
        
        onMounted(() => {
            updatePage(PAGE.UPCOMING);
        });

        return {
        }
    },
    template: `
    <div>
        <h1>Upcoming</h1>
    </div>
    `
}
export { Upcoming };

