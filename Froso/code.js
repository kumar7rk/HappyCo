const puppeteer = require('puppeteer');

(async () => {
	
	//const browser = await puppeteer.launch({headless: false}); // default is true

  const browser = await puppeteer.launch();
  const page = await browser.newPage();

  await page.goto('https://vequity.appfolio.com/connect/users/sign_in');
  
  const name = await page.evaluate(() => document.querySelector('.footer__info > p').innerText)
  const phone = await page.evaluate(() => document.querySelector('.footer__info > a').innerText)
  const website = await page.evaluate(() => document.querySelector('.footer__info > p > a ').href)
	
	console.log(name+ "\n" + phone+ "\n" + website)

  await browser.close();
})();