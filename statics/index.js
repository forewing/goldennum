var app = new Vue({
    el: '#app',
    data: {

    },
    methods: {
        updateDashboardRoomId(id) {
            if (id == null) {
                id = parseInt(this.$refs.tempRoomIdInput.value)
            }
            if (id > 0) {
                localStorage.setItem(KEY_SAVED_ROOM_ID, id);
                this.$refs.panelDashboard.updateRoomId(id);
            }
        },
    }
})