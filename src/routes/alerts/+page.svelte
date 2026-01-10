<script>
    import {alerts, alertsError, getAlerts, postAlerts, putAlert, deleteAlert} from "$lib/stores/alerts.js";
    import DeleteModal from "../../components/delete-modal.svelte"
    import {onMount} from "svelte";
    import {SvelteMap} from "svelte/reactivity";
    import {getAgenda} from "$lib/stores/agendas.js";

    let newEmail = $state("")
    let newAgendaId = $state("")
    let editable = new SvelteMap()
    let alertAgenda = new SvelteMap()

    let deleteModal = $state();

    onMount(() => {
        getAlerts()
    })

    alerts.subscribe((value) => {
        value.forEach((a) => {
            if(!alertAgenda.get(a.agendaId)){
                getAgenda(a.agendaId)
                    .then((data) => {
                        alertAgenda.set(a.agendaId, data)
                    })
                    .catch(()=>{})
            }
        })
    })


    function addAlert() {
        console.log(`Adding alert with email ${newEmail} and agendaId ${newAgendaId}`)
        postAlerts(newEmail, newAgendaId)
            .then(() => {
                newEmail = ""
                newAgendaId = ""
            })
            .catch(() => {
            })
    }

    function deleteAlertById(id) {
        console.log(`Deleting alert ${id}`)
        deleteAlert(id)
            .then(() => {
                editable.set(id, false)
            })
            .catch(() => {
            })
    }

    function editAlert(id, email, agendaId) {
        console.log(`Editing alert ${id} with email ${email} and agendaId ${agendaId}`)
        putAlert(id, email, agendaId)
            .then(() => {
                editable.set(id, false)
            })
            .catch(() => {
            })
    }
</script>

<div class="d-flex flex-column">
    <h1 class="text-center color-yellow mb-lg-5">Alerts</h1>
    <div class="d-flex flex-column mb-5">
        <h5 class="text-center fw-semibold mb-2">Add an alert : </h5>
        <form class="d-flex justify-content-center" onsubmit="{addAlert}">
            <div class="d-flex me-4">
                <label for="email" class="fs-5 align-self-center me-2">UCA Email</label>
                <input id="email" class="form-control-sm" placeholder="justine.bachelard@ext.uca.fr" bind:value={newEmail}/>
            </div>
            <div class="d-flex me-4">
                <label for="agenda" class="fs-5 align-self-center me-2">Agenda ID</label>
                <input id="agenda" class="form-control-sm" placeholder="6ce750e5-05b8-4c4c-8fee-d9e381dbf364 " bind:value={newAgendaId}/>
            </div>
            <input class="btn btn-yellow" type="submit" value="Submit"/>
        </form>
    </div>
    {#if $alertsError}
        <div class="text-center text-danger form-text fs-5">
            An error occurred : {$alertsError}
        </div>
    {/if}
    <table class="table w-50 align-self-center">
        <thead>
        <tr>
            <th class="color-yellow" scope="col">Email</th>
            <th class="color-yellow" scope="col">Agenda</th>
            <th class="color-yellow" scope="col">ID</th>
            <th scope="col">Actions</th>
        </tr>
        </thead>
        <tbody>
        {#each $alerts as alert}
            <tr>
                <td>
                    {#if !editable.get(alert.id)}
                        <p class="color-yellow fw-bold">{alert.email}</p>
                    {:else}
                        <input class="form-control-sm" bind:value={alert.email} />
                    {/if}
                </td>
                <td>
                    {#if !editable.get(alert.id)}
                        {alert.agendaId} <br/>
                        {alertAgenda.get(alert.agendaId)?.name ? alertAgenda.get(alert.agendaId)?.name : "Pb while getting agenda name"}
                    {:else}
                        <input class="form-control-sm" bind:value={alert.agendaId} />
                    {/if}
                </td>
                <td>{alert.id}</td>
                <td class="justify-content-around">
                    {#if !editable.get(alert.id)}
                        <button class="bg-transparent border-0" onclick="{() => {editable.set(alert.id, true)}}">
                            <span class="fa fa-pen-to-square text-warning"></span>
                        </button>
                        <button class="bg-transparent border-0" onclick={() => { deleteModal.show(alert.id, alert.id) }}>
                            <span class="fa fa-trash-can text-danger"></span>
                        </button>
                    {:else}
                        <button class="bg-transparent border-0" onclick="{() => {editAlert(alert.id, alert.email, alert.agendaId)}}">
                            <span class="fa fa-check" style="color: #0b6f33"></span>
                        </button>
                    {/if}
                </td>
            </tr>
        {/each}
        </tbody>
    </table>
    <DeleteModal bind:this={deleteModal} uniqueName="alerts" deleteFunction="{deleteAlertById}" />
</div>