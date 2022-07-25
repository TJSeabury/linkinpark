<script lang="ts">
	import { target } from '../stores';

	let targetValue: string;
	target.subscribe((v) => {
		targetValue = v;
	});

	const handleInput = (ev: Event) => {
		let value = (ev as InputEvent).data || '';
		target.set(value);
	};

	export let handleSubmit: (ev: Event) => void;
</script>

<form id="crawl" on:submit|preventDefault={handleSubmit}>
	<div>
		<input
			id="url"
			type="text"
			name="url"
			placeholder="https://example.com/"
			on:input={handleInput}
		/>
	</div>
	<div>
		<input type="submit" value="Crawl" />
	</div>
</form>
