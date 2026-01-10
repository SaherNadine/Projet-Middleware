<script>
    import {goto} from "$app/navigation";
</script>

<nav class="navbar bg-secondary align-content-start" style="width: fit-content">
    <ul class="d-flex flex-column navbar-nav text-center">
        <li class="nav-item mb-4 clickable-item">
            <a class="clickable-item" href="/agendas">
                <span class="h1 fas fa-calendar-plus"></span>
            </a>
        </li>
        <li class="nav-item mb-4 clickable-item">
            <a href="/alerts">
                <span class="h1 fas fa-bell"></span>
            </a>
        </li>
        <li class="nav-item mb-4 clickable-item">
            <a class="clickable-item" href="/events">
                <span class="h1 fas fa-calendar-days"></span>
            </a>
        </li>
    </ul>
</nav>

<style>
    .clickable-item:hover {
        background-color: rgb(var(--bs-secondary-bg-rgb), 0.2);
        cursor: pointer;
    }

    .nav-item {
        padding: 15px 15px 10px 15px;
    }

    a {
        color: white;
        text-decoration: none;
    }
</style>