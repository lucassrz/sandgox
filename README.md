# Sandgox

## Problématique
Créer un jeu sandbox avec différents élements avec une physique.

## Le projet

### Micro optimisation
- Mémoïsation des objet rectangle de largeur n.
- On réutilise au maximum les objets créés.
- Les fonctions gérants la physique ne sont pas instanciées par les cellules de la grille, de même pour la densitée. 

### Macro optimisation
- On récupère la frame précédente et on dessine par dessus les carrés qui sont en mouvement.
- On dessine des rectangle horizontalement si des élements de la même couleur sont collés.

### Benchmark
Comme le programme fonctionne avec une interface graphique,
nous avons du implémenter notre propre mode de manière de
benchmark le programme.

Pour cela, nous avons ajouter un flag falcultatif
"-benchmark true" au lancement du programme. il permet de
rentrer dans ce mode et la fenêtre ce fermera lorsque le
programme aura passé la 100ème frame.

C'est notre valeur de mesure pour ce benchmark.

Afin d'être le plus proche possible de l'utilisation,
Le benchmark est réalisé avec tout les élements
(sable, eau, métal, générateur d'eau, trou noir)
ainsi que les différent état possible (en mouvement, statique)

### Résultats

**V1**: **xx.xx** secondes pour 100 frames.

**Programme final**: **xx.xx** secondes pour 100 frames.