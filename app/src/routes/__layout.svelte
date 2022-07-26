<script>
	import '../app.css';

	const _app_version = '0.3.2';
</script>

<slot />

<footer class="flex justify-between p-8">
	<small>{_app_version}</small>
	<small>&copy; Tyler Seabury 2022</small>
</footer>
