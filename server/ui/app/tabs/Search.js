import { currentPage, PAGE, theme, updatePage } from "../utils.js";
import { onMounted, ref } from "../vue.js";


const Search = {
    props: {

    },
    components: {

    },
    setup: (props) => {

        onMounted(() => {
            updatePage(PAGE.SEARCH);
        });


        return {
        }
    },
    template: `
    <div>
        <h1>Search</h1>
    </div>
    `
}
export { Search };

