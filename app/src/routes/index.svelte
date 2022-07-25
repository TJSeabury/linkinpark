<script lang="ts">
	import { target } from '../stores';
	import CrawlForm from '$lib/CrawlForm.svelte';
	import { wait, diff } from '$lib/Utils';
	import { apiHost } from '$lib/env';

	let targetValue: string;

	target.subscribe((v) => {
		targetValue = v;
	});

	let id: string = '';
	let status: string = 'idle';
	let elapsedTime: number = 0;
	let linksFound: number = 0;
	let linksCrawled: number = 0;
	let done: boolean = false;
	let downloadLink: string = '';

	async function handleSubmit(ev: Event) {
		ev.preventDefault();
		console.log('Yee ha');
		const response = await fetch(`${apiHost}/api/start/`, {
			method: 'POST',
			mode: 'cors',
			cache: 'no-cache',
			headers: {
				'Content-Type': 'application/json'
			},
			body: JSON.stringify({
				uuid: null,
				message: targetValue || ''
			})
		});

		let rData = await response.json();

		if (response.status === 422) {
			status = rData.message;
			return;
		}

		const uuid = rData.uuid;
		status = rData.message;
		done = false;

		while (!done) {
			const checkResponse = await fetch(`${apiHost}/api/check/`, {
				method: 'POST',
				mode: 'cors',
				cache: 'no-cache',
				credentials: 'same-origin',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify({
					uuid: uuid,
					message: ''
				})
			});

			rData = await checkResponse.json();

			if (rData.Status === 'done') done = true;

			id = rData.Uuid;
			status = rData.Status;
			elapsedTime = diff() / 1000;
			linksFound = rData.LinksFound;
			linksCrawled = rData.LinksCrawled;

			wait(100);

			if (done) {
				downloadLink = `${apiHost}/api/finish/${uuid}`;
			}
		}
	}

	function reset() {
		id = '';
		status = 'idle';
		elapsedTime = 0;
		linksFound = 0;
		linksCrawled = 0;
		done = false;
		downloadLink = '';
	}
</script>

<svelte:head>
	<title>Linkinpark</title>
</svelte:head>

<main>
	<header>
		<h1>Linkinpark</h1>
		<p>Crawling in my crawl.</p>
	</header>
	<section>
		<CrawlForm {handleSubmit} />
		<div id="crawl-status">
			<p>ID: {id}</p>
			<p>Status: {status}</p>
			<p>Links found: {linksFound}</p>
			<p>Links crawled: {linksCrawled}</p>
			<p>Elapsed: {elapsedTime} seconds</p>
		</div>
		<div id="download-link">
			<p>
				{#if done}
					<a href={downloadLink} on:click={reset}>Download Report</a>
				{/if}
			</p>
		</div>
	</section>
</main>
