<svelte:head>
	<title>Linkinpark</title>
</svelte:head>

<main>
	<header>
		<h1>Linkinpark</h1>
		<p>Crawling in my crawl.</p>
	</header>
	<section>
		<form id="crawl" action="/api/start/" method="post">
			<div>
				<input id="url" type="text" name="url" placeholder="https://horacesmithfund.org/" />
			</div>
			<div>
				<input type="submit" value="Crawl" />
			</div>
		</form>
		<div id="status">
			<p />
		</div>
		<div id="download-link">
			<p />
		</div>
	</section>
</main>
