<script>
    let {
        uniqueName,
        deleteFunction,
    } = $props()

    let modalName = $derived(`deleteModal-${uniqueName}`)

    export function show(id, name) {
        if (deleteModal === undefined) {
            deleteModal = new bootstrap.Modal(`#${modalName}`);
        }
        itemId = id;
        itemName = name;
        deleteModal.show();
    }


    let itemId = $state("");
    let itemName = $state("");
    let deleteModal = $state();

    function deleteResource() {
        deleteFunction(itemId);
        deleteModal.hide();
    }

    function reset() {
        deleteModal.hide();
    }
</script>

<div class="modal fade modal-lg" id="{modalName}" aria-hidden="true" data-bs-backdrop="static">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h1 class="modal-title fs-5">
                    Deletion confirmation
                </h1>
                <button type="button" class="btn-close" data-bs-dismiss="modal" onclick={reset}></button>
            </div>

            <div class="modal-body d-flex align-items-start justify-content-between flex-md-column">
                Are you sure you want to delete {itemName} ?
            </div>

            <div class="modal-footer">
                <button type="button" class="btn btn-secondary rounded-0" data-bs-dismiss="modal" onclick={reset}>
                    Cancel
                </button>
                <button type="button" class="btn btn-danger delete-button rounded-0 bg-danger" onclick={deleteResource}>
                    Delete
                </button>
            </div>
        </div>
    </div>
</div>