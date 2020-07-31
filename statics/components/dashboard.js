Vue.component('dashboard', {
    props: [],
    template: "#dashboard-component",
    data: function () {
        return {
            data: "{}",
        }
    },
    methods: {
        refreshRoom() {
            localStorage.setItem(KEY_ROOM_ID, 1);
            const roomId = localStorage.getItem(KEY_ROOM_ID);
            if (!roomId) {
                return;
            }
            fetch(URL_ROOM_INFO + roomId)
                .then(response => response.json())
                .then(data => {
                    this.data = data;
                });
        }
    }
})
