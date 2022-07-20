function wait ( ms ) {
  const start = Date.now();
  let now = start;
  while ( now - start < ms ) {
    now = Date.now();
  }
}

const data = {
  target: ''
};
const crawlForm = document.querySelector( '#crawl' );

const status = document.querySelector( '#status > p' );
const setStatus = value => {
  status.innerHTML = value;
};

const downloadLink = document.querySelector( '#download-link > p' );
const clearDownloadLink = () => {
  downloadLink.innerHTML = null;
};
const setDownloadLink = linkElement => {
  clearDownloadLink();
  downloadLink.appendChild( linkElement );
};

crawlForm.addEventListener( 'input', ev => {
  if ( ev.target.id === 'url' ) {
    data.target = ev.target.value;
    console.log( data.target );
  }
} );

const handleSubmit = async ev => {
  ev.preventDefault();

  const response = await fetch(
    crawlForm.getAttribute( 'action' ),
    {
      method: 'POST',
      mode: 'cors',
      cache: 'no-cache',
      credentials: 'same-origin',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify( {
        uuid: null,
        message: data.target || ''
      } )
    }
  );

  let rData = await response.json();

  if ( response.status === 422 ) {
    setStatus( rData.message );
    return;
  }

  const uuid = rData.uuid;

  setStatus( rData.message );

  const diff = ( () => {
    const start = Date.now();
    return () => Date.now() - start;
  } )();
  const limit = ms => diff() < ms;

  let done = false;

  while ( rData.Status !== 'done' ) {
    if ( limit( 1000 * 600 ) ) {
      const checkResponse = await fetch(
        '/api/check/',
        {
          method: 'POST',
          mode: 'cors',
          cache: 'no-cache',
          credentials: 'same-origin',
          headers: {
            'Content-Type': 'application/json'
          },
          body: JSON.stringify( {
            uuid: uuid,
            message: ''
          } )
        } );

      rData = await checkResponse.json();

      setStatus( `${rData.Uuid} | ${rData.Status} | Elapsed: ${diff() / 1000} seconds | LinksFound: ${rData.LinksFound} | LinksCrawled: ${rData.LinksCrawled}` );

      if ( rData.Status === 'done' ) done = true;

      wait( 100 );

    } else {
      setStatus( `${rData.Uuid} | Timed out!! | Elapsed: ${diff() / 1000} seconds | LinksFound: ${rData.LinksFound} | LinksCrawled: ${rData.LinksCrawled}` );
      break;
    }

  }

  if ( done ) {
    let dlLink = document.createElement( 'a' );
    dlLink.setAttribute(
      'href',
      `/api/finish/${uuid}`
    );
    dlLink.innerHTML = 'Dowload Report';
    dlLink.addEventListener(
      'click',
      clearDownloadLink
    );
    setDownloadLink( dlLink );
  }

};

crawlForm.addEventListener( 'submit', handleSubmit );