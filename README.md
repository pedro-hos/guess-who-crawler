# guess-who-crawler

This Golang-based web crawler is designed to scrape Wikipedia pages to retrieve a list of people born in Brazil, categorized by state and city.

## The Who-Guess Game
The Who-Guess game is designed to be similar to a profile-based guessing game. In this game, a person is selected randomly, and players must guess who it is based on a series of clues provided one by one. The goal is to identify the selected person with as few clues as possible. You can choose individuals from any state or city in Brazil.

## Guess-Who Crawler

This Golang-based web crawler is designed to scrape Wikipedia pages to retrieve a list of people born in Brazil, categorized by state and city. So, to run you need to install and have Golang installed and running.

This is the pages scrapped:

1. First we scrap the Naturals (born in) Brazil by Federated Unit (State) [Categoria:Naturais do Brasil por unidade federativa](https://pt.wikipedia.org/wiki/Categoria:Naturais_do_Brasil_por_unidade_federativa);
2. Second we go on each page and get the Cites present on each state. For example [Naturals from São Paulo State](https://pt.wikipedia.org/w/index.php?title=Categoria:Naturais_do_estado_de_S%C3%A3o_Paulo&subcatuntil=Jacare%C3%AD%0ANaturais+de+Jacare%C3%AD#mw-subcategories);
3. Finally, we go on each people born on each state city, for example [born in São José dos Campos](https://pt.wikipedia.org/wiki/Categoria:Naturais_de_S%C3%A3o_Jos%C3%A9_dos_Campos) page and save it;

After all this steps, we get trainned LLM in order to provide the clues. (**In progress**)

## Todo
This is the todo list in order to get the Beta version working

- [x] Create and configure the initial Golang project;
- [x] Use [Colly](https://github.com/gocolly/colly) as libary to scrap the data for all pages;
- [x] Analyze and scrap the Naturals (born in) Brazil by Federated Unit (State) [Categoria:Naturais do Brasil por unidade federativa](https://pt.wikipedia.org/wiki/Categoria:Naturais_do_Brasil_por_unidade_federativa);
- [x] Save all states scrapped on the database;
- [x] Iterate over the states scrapped before and get all cities belong to the respective state. For example [Naturals from São Paulo State](https://pt.wikipedia.org/w/index.php?title=Categoria:Naturais_do_estado_de_S%C3%A3o_Paulo&subcatuntil=Jacare%C3%AD%0ANaturais+de+Jacare%C3%AD#mw-subcategories) and save it;
- [x] Iterate over the cities scrapped before, for example, [born in São José dos Campos](https://pt.wikipedia.org/wiki/Categoria:Naturais_de_S%C3%A3o_Jos%C3%A9_dos_Campos) page, and save the person name as the card answer and the wikipedia link, for example: [Cassinano Ricardo](https://pt.wikipedia.org/wiki/Cassiano_Ricardo) page save those information.
- [ ] Scrap the Avatar and save the link
- [ ] **To be decided** - We need to use some local LLM with [Podman AI](https://podman-desktop.io/docs/ai-lab) or [InstructLab AI](https://instructlab.ai/) in order to training a Model that will be able to:
    - [ ] Read the Wikipedia saved before and provide the Clues based on the respective people page OR
    - [ ] Read the output from the page Wikipedia API for example: `https://pt.wikipedia.org/w/api.php?action=query&prop=revisions&rvprop=content&format=json&titles=Aramis%20Brito&rvsection=0` and provided the clues based on that response

## Other repositories for the Guess-Who project

- (**In progress ...**) [Guess-Who API](https://github.com/pedro-hos/guess-who-api) - This is the API project for all data scrapped with the Guess-Who Crawler page;
- (**In progress ...**) [Guess-Who Web](https://github.com/pedro-hos/guess-who-web) - This is the front-end page, the game page.