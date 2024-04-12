#  Dossier de synthèse CESI

Je pense que tout le monde sera d'accord pour dire que les dossiers de synthèse ne sont pas très beaux.
Et même au-delà, les calculs de note peuvent être fastidieux. Ce projet offre donc une belle interface
permettant de visualiser son dossier, et de changer ses notes pour voir si ça passe !

https://github.com/julien-wff/cesi-dossier-synthese/assets/50249422/31696820-5886-430d-bb2e-a335c372b045

## Fonctionnement

Pour faire simple, le PDF est envoyé au serveur, qui va en extraire les tableaux à l'aide de 
[Tabula](https://github.com/tabulapdf/tabula-java).
Ces tableaux vont ensuite être parsés, et renvoyés au client.

## Confidentialité

Les PDF ne sont jamais conservés. Une fois traités, ils sont directement supprimés du serveur.
C'est également le cas s'il y a une erreur, supprimés !
Si vous n'avez pas confiance, le projet est opensource, hostez-le vous même !

## Affiliations

Ce projet n'est pas associé à [CESI](https://www.cesi.fr), mais l'est au 
[BDE du CESI Nancy](https://bdecesinancy.fr) qui propose une [version hébergée](https://dossier.bdecesinancy.fr)
utilisable gratuitement.
