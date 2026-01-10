<script>
    import Navbar from "../components/navbar.svelte";
    import SideBar from "../components/sidebar.svelte";
</script>

<div class="d-flex flex-column">
    <Navbar></Navbar>
    <div class="d-flex flex-row" style="min-height: 95vh">
        <SideBar></SideBar>
        <div  class="mx-auto" style="margin-top: 4rem; padding: 2rem;width: 90%">
            <slot></slot>
        </div>
    </div>
</div>
