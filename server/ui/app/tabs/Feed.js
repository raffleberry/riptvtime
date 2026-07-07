import { currentPage, PAGE, theme, updatePage } from "../utils.js";
import { onMounted, ref } from "../vue.js";


const Feed = {
    props: {

    },
    components: {

    },
    setup: (props) => {

        onMounted(() => {
            updatePage(PAGE.FEED);
        });


        return {
        }
    },
    template: `
    <div>
        <h1>Feed</h1>
    </div>
    `
}
export { Feed };

