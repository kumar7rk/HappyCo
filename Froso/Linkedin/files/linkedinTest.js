const puppeteer = require('puppeteer');
const CREDS = require('./creds');

(async () => {
	const browser = await puppeteer.launch({
		headless: false
	});

	const page = await browser.newPage();
	
	await page.goto("https://www.linkedin.com");

	const USERNAME_SELECTOR = '#login-email';
  	const PASSWORD_SELECTOR = '#login-password';
  	const BUTTON_SELECTOR = '#login-submit';

  	await page.click(USERNAME_SELECTOR);
  	await page.keyboard.type(CREDS.username);

  	await page.click(PASSWORD_SELECTOR);
  	await page.keyboard.type(CREDS.password);
	
	await page.click(BUTTON_SELECTOR);
    await page.waitForNavigation();
  	
	var url = "https://www.linkedin.com/in/robert-wise/";

	await page.goto(url);
	
	await page.waitFor(2 * 1000);
	await page.evaluate(_ => {
	  window.scrollBy(0, window.innerHeight);
	});
	await page.waitFor(2 * 1000);
   const name = await page.evaluate(() => document.querySelector
    ('div.pv-top-card-v2-section__info.mr5 > div.display-flex.align-items-center > h1').textContent)
   console.log(name);
   
   const title = await page.evaluate(() => document.querySelector
   	('div.pv-entity__summary-info.pv-entity__summary-info--v2 >h3').textContent)
   
   console.log("title "+title);
   
   const companyNameExtra = await page.evaluate(() => document.querySelector
   	('#experience-section > ul > div > li > a > div > h4').textContent)
  /* ul
	div
	li
	a
	div
	h3*/
   console.log("companyNameExtra "+companyNameExtra);
   
   const companyName = await page.evaluate(() => document.querySelector
   	('div > h4 > span:nth-child(2)').textContent)
   console.log("companyName "+companyName);
   
   //#ember6250 > div.pv-entity__summary-info.pv-entity__summary-info--v2 > h4.Sans-17px-black-85\25 > span.pv-entity__secondary-title
   //#ember6250 > div.pv-entity__summary-info.pv-entity__summary-info--v2 > h4.pv-entity__date-range.inline-block.Sans-15px-black-70\25 > span:nth-child(2)
   
   const jobDates = await page.evaluate(() => document.querySelector
   	('div > h4.pv-entity__date-range.inline-block > span:nth-child(2)').textContent)
   console.log("jobDates "+jobDates);
   //#ember6250 > div.pv-entity__summary-info.pv-entity__summary-info--v2 > h4:nth-child(4) > span.pv-entity__bullet-item-v2

  const jobDuration = await page.evaluate(() => document.querySelector
   	('div > h4:nth-child(4)	 > span.pv-entity__bullet-item-v2').textContent)
   console.log("jobDuration "+jobDuration);
 await browser.close();
})();